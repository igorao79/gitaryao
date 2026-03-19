package archive

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Manager handles backup and restore of bare git repos to/from the database.
type Manager struct {
	DB       *sql.DB
	ReposDir string
}

// NewManager creates a new archive manager.
func NewManager(db *sql.DB, reposDir string) *Manager {
	return &Manager{DB: db, ReposDir: reposDir}
}

// BackupRepo archives a bare repo and stores it in the database.
func (m *Manager) BackupRepo(repoID int64, owner, name string) error {
	repoPath := filepath.Join(m.ReposDir, owner, name+".git")

	data, err := tarGzDir(repoPath)
	if err != nil {
		return fmt.Errorf("archive repo %s/%s: %w", owner, name, err)
	}

	_, err = m.DB.Exec(
		`INSERT INTO repo_archives (repo_id, data, updated_at) VALUES (?, ?, ?)
		 ON CONFLICT(repo_id) DO UPDATE SET data = excluded.data, updated_at = excluded.updated_at`,
		repoID, data, time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("save archive %s/%s: %w", owner, name, err)
	}

	fmt.Printf("Backed up repo %s/%s (%d bytes)\n", owner, name, len(data))
	return nil
}

// BackupByOwnerAndName looks up repo_id and backs up the repo.
func (m *Manager) BackupByOwnerAndName(owner, name string) error {
	var repoID int64
	err := m.DB.QueryRow(
		`SELECT r.id FROM repositories r JOIN users u ON r.owner_id = u.id
		 WHERE u.username = ? AND r.name = ?`, owner, name,
	).Scan(&repoID)
	if err != nil {
		return fmt.Errorf("find repo %s/%s: %w", owner, name, err)
	}
	return m.BackupRepo(repoID, owner, name)
}

// RestoreAll restores all archived repos from the database to disk.
func (m *Manager) RestoreAll() error {
	rows, err := m.DB.Query(
		`SELECT ra.data, u.username, r.name
		 FROM repo_archives ra
		 JOIN repositories r ON ra.repo_id = r.id
		 JOIN users u ON r.owner_id = u.id`,
	)
	if err != nil {
		return fmt.Errorf("query archives: %w", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var data []byte
		var owner, name string
		if err := rows.Scan(&data, &owner, &name); err != nil {
			fmt.Printf("Warning: skip archive row: %v\n", err)
			continue
		}

		repoPath := filepath.Join(m.ReposDir, owner, name+".git")

		// Skip if repo already exists on disk
		if _, err := os.Stat(repoPath); err == nil {
			continue
		}

		if err := os.MkdirAll(filepath.Dir(repoPath), 0755); err != nil {
			fmt.Printf("Warning: mkdir %s: %v\n", repoPath, err)
			continue
		}

		if err := untarGz(data, repoPath); err != nil {
			fmt.Printf("Warning: restore %s/%s: %v\n", owner, name, err)
			continue
		}

		fmt.Printf("Restored repo %s/%s (%d bytes)\n", owner, name, len(data))
		count++
	}

	if count > 0 {
		fmt.Printf("Restored %d repos from database\n", count)
	}
	return rows.Err()
}

// tarGzDir creates a tar.gz archive of a directory.
func tarGzDir(srcDir string) ([]byte, error) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	baseDir := filepath.Base(srcDir)

	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path from parent of srcDir
		relPath, err := filepath.Rel(filepath.Dir(srcDir), path)
		if err != nil {
			return err
		}
		// Normalize to forward slashes
		relPath = strings.ReplaceAll(relPath, "\\", "/")

		// Skip the base directory entry itself if it matches
		if relPath == baseDir && info.IsDir() {
			// Still need to add directory entry
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = relPath

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(tw, f)
		return err
	})
	if err != nil {
		return nil, err
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}
	if err := gw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// untarGz extracts a tar.gz archive into destDir's parent, creating the repo directory.
func untarGz(data []byte, destDir string) error {
	gr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	parentDir := filepath.Dir(destDir)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Sanitize path
		cleanName := filepath.Clean(header.Name)
		if strings.Contains(cleanName, "..") {
			continue
		}

		target := filepath.Join(parentDir, cleanName)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			f.Close()
		}
	}

	return nil
}
