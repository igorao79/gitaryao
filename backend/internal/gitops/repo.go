package gitops

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	git "github.com/go-git/go-git/v5"
)

// GetRepoPath returns the filesystem path for a bare repository.
func GetRepoPath(basePath, owner, name string) string {
	return filepath.Join(basePath, owner, name+".git")
}

// InitBareRepo creates a new bare git repository on disk.
func InitBareRepo(basePath, owner, name string) error {
	repoPath := GetRepoPath(basePath, owner, name)

	if err := os.MkdirAll(filepath.Dir(repoPath), 0755); err != nil {
		return fmt.Errorf("create owner dir: %w", err)
	}

	_, err := git.PlainInit(repoPath, true)
	if err != nil {
		return fmt.Errorf("init bare repo: %w", err)
	}

	return nil
}

// RepoExists checks if a bare repository exists on disk.
func RepoExists(basePath, owner, name string) bool {
	repoPath := GetRepoPath(basePath, owner, name)
	info, err := os.Stat(repoPath)
	return err == nil && info.IsDir()
}

// ListRepos scans the repos directory and returns all repos grouped by owner.
func ListRepos(basePath string) ([]RepoInfo, error) {
	var repos []RepoInfo

	owners, err := os.ReadDir(basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return repos, nil
		}
		return nil, err
	}

	for _, ownerEntry := range owners {
		if !ownerEntry.IsDir() {
			continue
		}
		owner := ownerEntry.Name()

		repoEntries, err := os.ReadDir(filepath.Join(basePath, owner))
		if err != nil {
			continue
		}

		for _, repoEntry := range repoEntries {
			if !repoEntry.IsDir() || !strings.HasSuffix(repoEntry.Name(), ".git") {
				continue
			}
			name := strings.TrimSuffix(repoEntry.Name(), ".git")
			repos = append(repos, RepoInfo{
				Owner: owner,
				Name:  name,
			})
		}
	}

	return repos, nil
}

// RepoInfo holds basic repository metadata.
type RepoInfo struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}
