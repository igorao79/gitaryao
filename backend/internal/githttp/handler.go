package githttp

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Handler implements Git Smart HTTP protocol.
// It delegates to git-upload-pack and git-receive-pack binaries.
type Handler struct {
	ReposDir string
}

// NewHandler creates a new Git Smart HTTP handler.
func NewHandler(reposDir string) *Handler {
	return &Handler{ReposDir: reposDir}
}

// InfoRefs handles GET /{owner}/{repo}.git/info/refs?service=...
func (h *Handler) InfoRefs(w http.ResponseWriter, r *http.Request) {
	repoPath, err := h.resolveRepo(r)
	if err != nil {
		http.Error(w, "Repository not found", http.StatusNotFound)
		return
	}

	service := r.URL.Query().Get("service")
	if service != "git-upload-pack" && service != "git-receive-pack" {
		http.Error(w, "Invalid service", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", fmt.Sprintf("application/x-%s-advertisement", service))
	w.Header().Set("Cache-Control", "no-cache")

	// Write service announcement pkt-line
	serverAdvert := fmt.Sprintf("# service=%s\n", service)
	w.Write(pktLine(serverAdvert))
	w.Write(pktFlush())

	// Run git command to get refs
	cmd := exec.Command("git", service[4:], "--stateless-rpc", "--advertise-refs", repoPath)
	cmd.Env = append(os.Environ(), "GIT_PROTOCOL=version=2")

	out, err := cmd.Output()
	if err != nil {
		http.Error(w, "Failed to get refs", http.StatusInternalServerError)
		return
	}

	w.Write(out)
}

// UploadPack handles POST /{owner}/{repo}.git/git-upload-pack (clone/fetch)
func (h *Handler) UploadPack(w http.ResponseWriter, r *http.Request) {
	h.serviceRPC(w, r, "upload-pack")
}

// ReceivePack handles POST /{owner}/{repo}.git/git-receive-pack (push)
func (h *Handler) ReceivePack(w http.ResponseWriter, r *http.Request) {
	h.serviceRPC(w, r, "receive-pack")
}

func (h *Handler) serviceRPC(w http.ResponseWriter, r *http.Request, service string) {
	repoPath, err := h.resolveRepo(r)
	if err != nil {
		http.Error(w, "Repository not found", http.StatusNotFound)
		return
	}

	// Verify content type
	expectedCT := fmt.Sprintf("application/x-git-%s-request", service)
	if r.Header.Get("Content-Type") != expectedCT {
		http.Error(w, "Invalid content type", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", fmt.Sprintf("application/x-git-%s-result", service))
	w.Header().Set("Cache-Control", "no-cache")

	// Handle gzip-encoded request body
	var body io.Reader = r.Body
	if r.Header.Get("Content-Encoding") == "gzip" {
		body, err = gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, "Failed to decompress", http.StatusInternalServerError)
			return
		}
	}

	cmd := exec.Command("git", service, "--stateless-rpc", repoPath)
	cmd.Env = append(os.Environ(), "GIT_PROTOCOL=version=2")
	cmd.Stdin = body
	cmd.Stdout = w

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// Log stderr for debugging but don't expose to client
		fmt.Fprintf(os.Stderr, "git %s error: %s\n", service, stderr.String())
		return
	}
}

// resolveRepo extracts owner/repo from the URL and returns the filesystem path.
// Expects URL pattern: /{owner}/{repo}.git/...
func (h *Handler) resolveRepo(r *http.Request) (string, error) {
	// Extract owner and repo from chi URL params or parse path
	path := r.URL.Path
	parts := strings.Split(strings.Trim(path, "/"), "/")

	if len(parts) < 2 {
		return "", fmt.Errorf("invalid path")
	}

	owner := parts[0]
	repoWithGit := parts[1]

	// Validate owner and repo names (prevent path traversal)
	if strings.Contains(owner, "..") || strings.Contains(repoWithGit, "..") {
		return "", fmt.Errorf("invalid path")
	}
	if strings.Contains(owner, "/") || strings.Contains(owner, "\\") {
		return "", fmt.Errorf("invalid owner")
	}

	repoPath := filepath.Join(h.ReposDir, owner, repoWithGit)

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return "", fmt.Errorf("repo not found: %s/%s", owner, repoWithGit)
	}

	return repoPath, nil
}

// pktLine encodes a string in git pkt-line format.
// Format: 4-byte hex length (including length bytes) + data
func pktLine(data string) []byte {
	length := len(data) + 4
	return []byte(fmt.Sprintf("%04x%s", length, data))
}

// pktFlush returns a flush packet (0000).
func pktFlush() []byte {
	return []byte("0000")
}
