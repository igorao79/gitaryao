package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"gitserv/internal/gitops"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-chi/chi/v5"
)

// TreeEntry represents a file or directory in a git tree.
type TreeEntry struct {
	Name string `json:"name"`
	Type string `json:"type"` // "blob" or "tree"
	Size int64  `json:"size"`
}

// CommitInfo represents a commit in the log.
type CommitInfo struct {
	Hash    string    `json:"hash"`
	Message string    `json:"message"`
	Author  string    `json:"author"`
	Email   string    `json:"email"`
	Date    time.Time `json:"date"`
}

// BranchInfo represents a branch reference.
type BranchInfo struct {
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
}

// openRepo is a helper to open a bare repo by owner/name.
func (s *Server) openRepo(owner, name string) (*git.Repository, error) {
	repoPath := gitops.GetRepoPath(s.ReposDir, owner, name)
	return git.PlainOpen(repoPath)
}

// resolveRef resolves a ref string (branch name, tag, or hash) to a commit.
func resolveRef(repo *git.Repository, ref string) (*object.Commit, error) {
	// Try as a branch first
	hash, err := repo.ResolveRevision(plumbing.Revision(ref))
	if err != nil {
		// Try as a full hash
		h := plumbing.NewHash(ref)
		commit, err2 := repo.CommitObject(h)
		if err2 != nil {
			return nil, err
		}
		return commit, nil
	}
	return repo.CommitObject(*hash)
}

// extractSubpath extracts the wildcard path portion after the ref segment.
// For routes like /repos/{owner}/{name}/tree/{ref}/* the wildcard captures
// everything after the ref. chi puts this in the "*" URL param.
func extractSubpath(r *http.Request) string {
	// chi wildcard param
	p := chi.URLParam(r, "*")
	return strings.TrimPrefix(p, "/")
}

// GetTree handles GET /api/repos/{owner}/{name}/tree/{ref} and
// GET /api/repos/{owner}/{name}/tree/{ref}/*
func (s *Server) GetTree(w http.ResponseWriter, r *http.Request) {
	owner := chi.URLParam(r, "owner")
	name := chi.URLParam(r, "name")
	ref := chi.URLParam(r, "ref")
	subpath := extractSubpath(r)

	repo, err := s.openRepo(owner, name)
	if err != nil {
		http.Error(w, "repository not found", http.StatusNotFound)
		return
	}

	commit, err := resolveRef(repo, ref)
	if err != nil {
		http.Error(w, "ref not found: "+err.Error(), http.StatusNotFound)
		return
	}

	rootTree, err := commit.Tree()
	if err != nil {
		http.Error(w, "failed to get tree: "+err.Error(), http.StatusInternalServerError)
		return
	}

	targetTree := rootTree
	if subpath != "" {
		entry, err := rootTree.FindEntry(subpath)
		if err != nil {
			http.Error(w, "path not found: "+subpath, http.StatusNotFound)
			return
		}
		if !entry.Mode.IsFile() {
			targetTree, err = repo.TreeObject(entry.Hash)
			if err != nil {
				http.Error(w, "failed to get subtree", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "path is a file, use blob endpoint", http.StatusBadRequest)
			return
		}
	}

	var entries []TreeEntry
	for _, e := range targetTree.Entries {
		entryType := "blob"
		var size int64
		if e.Mode.IsFile() {
			blob, err := repo.BlobObject(e.Hash)
			if err == nil {
				size = blob.Size
			}
		} else {
			entryType = "tree"
		}
		entries = append(entries, TreeEntry{
			Name: e.Name,
			Type: entryType,
			Size: size,
		})
	}

	if entries == nil {
		entries = []TreeEntry{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

// GetBlob handles GET /api/repos/{owner}/{name}/blob/{ref}/*
func (s *Server) GetBlob(w http.ResponseWriter, r *http.Request) {
	owner := chi.URLParam(r, "owner")
	name := chi.URLParam(r, "name")
	ref := chi.URLParam(r, "ref")
	filePath := extractSubpath(r)

	if filePath == "" {
		http.Error(w, "file path is required", http.StatusBadRequest)
		return
	}

	repo, err := s.openRepo(owner, name)
	if err != nil {
		http.Error(w, "repository not found", http.StatusNotFound)
		return
	}

	commit, err := resolveRef(repo, ref)
	if err != nil {
		http.Error(w, "ref not found: "+err.Error(), http.StatusNotFound)
		return
	}

	tree, err := commit.Tree()
	if err != nil {
		http.Error(w, "failed to get tree", http.StatusInternalServerError)
		return
	}

	file, err := tree.File(filePath)
	if err != nil {
		http.Error(w, "file not found: "+filePath, http.StatusNotFound)
		return
	}

	reader, err := file.Reader()
	if err != nil {
		http.Error(w, "failed to read file", http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	contents, err := io.ReadAll(reader)
	if err != nil {
		http.Error(w, "failed to read file contents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"path":    filePath,
		"size":    file.Size,
		"content": string(contents),
	})
}

// GetCommits handles GET /api/repos/{owner}/{name}/commits/{ref}
func (s *Server) GetCommits(w http.ResponseWriter, r *http.Request) {
	owner := chi.URLParam(r, "owner")
	name := chi.URLParam(r, "name")
	ref := chi.URLParam(r, "ref")

	repo, err := s.openRepo(owner, name)
	if err != nil {
		http.Error(w, "repository not found", http.StatusNotFound)
		return
	}

	commit, err := resolveRef(repo, ref)
	if err != nil {
		http.Error(w, "ref not found: "+err.Error(), http.StatusNotFound)
		return
	}

	logIter, err := repo.Log(&git.LogOptions{From: commit.Hash})
	if err != nil {
		http.Error(w, "failed to get log: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var commits []CommitInfo
	count := 0
	maxCommits := 50

	err = logIter.ForEach(func(c *object.Commit) error {
		if count >= maxCommits {
			return io.EOF
		}
		commits = append(commits, CommitInfo{
			Hash:    c.Hash.String(),
			Message: c.Message,
			Author:  c.Author.Name,
			Email:   c.Author.Email,
			Date:    c.Author.When,
		})
		count++
		return nil
	})
	// io.EOF is expected when we hit maxCommits or end of log
	if err != nil && err != io.EOF {
		http.Error(w, "failed to iterate commits", http.StatusInternalServerError)
		return
	}

	if commits == nil {
		commits = []CommitInfo{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commits)
}

// GetBranches handles GET /api/repos/{owner}/{name}/branches
func (s *Server) GetBranches(w http.ResponseWriter, r *http.Request) {
	owner := chi.URLParam(r, "owner")
	name := chi.URLParam(r, "name")

	repo, err := s.openRepo(owner, name)
	if err != nil {
		http.Error(w, "repository not found", http.StatusNotFound)
		return
	}

	// Determine the default branch by reading HEAD
	headRef, err := repo.Head()
	defaultBranch := ""
	if err == nil {
		defaultBranch = headRef.Name().Short()
	}

	branchIter, err := repo.Branches()
	if err != nil {
		http.Error(w, "failed to list branches", http.StatusInternalServerError)
		return
	}

	var branches []BranchInfo
	err = branchIter.ForEach(func(ref *plumbing.Reference) error {
		branchName := ref.Name().Short()
		branches = append(branches, BranchInfo{
			Name:      branchName,
			IsDefault: branchName == defaultBranch,
		})
		return nil
	})
	if err != nil {
		http.Error(w, "failed to iterate branches", http.StatusInternalServerError)
		return
	}

	if branches == nil {
		branches = []BranchInfo{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(branches)
}
