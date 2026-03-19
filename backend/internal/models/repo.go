package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Repository struct {
	ID            int64     `json:"id"`
	OwnerID       int64     `json:"owner_id"`
	OwnerName     string    `json:"owner_name,omitempty"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	IsPrivate     bool      `json:"is_private"`
	DefaultBranch string    `json:"default_branch"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type RepoStore struct {
	DB *sql.DB
}

// Create inserts a new repository record.
func (s *RepoStore) Create(ownerID int64, name, description string, isPrivate bool) (*Repository, error) {
	result, err := s.DB.Exec(
		`INSERT INTO repositories (owner_id, name, description, is_private) VALUES (?, ?, ?, ?)`,
		ownerID, name, description, isPrivate,
	)
	if err != nil {
		return nil, fmt.Errorf("insert repo: %w", err)
	}

	id, _ := result.LastInsertId()
	return s.FindByID(id)
}

// FindByID returns a repository by its database ID.
func (s *RepoStore) FindByID(id int64) (*Repository, error) {
	r := &Repository{}
	err := s.DB.QueryRow(
		`SELECT r.id, r.owner_id, u.username, r.name, r.description, r.is_private, r.default_branch, r.created_at, r.updated_at
		 FROM repositories r JOIN users u ON r.owner_id = u.id
		 WHERE r.id = ?`, id,
	).Scan(&r.ID, &r.OwnerID, &r.OwnerName, &r.Name, &r.Description, &r.IsPrivate, &r.DefaultBranch, &r.CreatedAt, &r.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find repo by id: %w", err)
	}
	return r, nil
}

// FindByOwnerAndName returns a repository by owner username and repo name.
func (s *RepoStore) FindByOwnerAndName(ownerName, repoName string) (*Repository, error) {
	r := &Repository{}
	err := s.DB.QueryRow(
		`SELECT r.id, r.owner_id, u.username, r.name, r.description, r.is_private, r.default_branch, r.created_at, r.updated_at
		 FROM repositories r JOIN users u ON r.owner_id = u.id
		 WHERE u.username = ? AND r.name = ?`, ownerName, repoName,
	).Scan(&r.ID, &r.OwnerID, &r.OwnerName, &r.Name, &r.Description, &r.IsPrivate, &r.DefaultBranch, &r.CreatedAt, &r.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find repo by owner/name: %w", err)
	}
	return r, nil
}

// ListByOwnerID returns all repositories for a user.
func (s *RepoStore) ListByOwnerID(ownerID int64) ([]Repository, error) {
	rows, err := s.DB.Query(
		`SELECT r.id, r.owner_id, u.username, r.name, r.description, r.is_private, r.default_branch, r.created_at, r.updated_at
		 FROM repositories r JOIN users u ON r.owner_id = u.id
		 WHERE r.owner_id = ? ORDER BY r.updated_at DESC`, ownerID,
	)
	if err != nil {
		return nil, fmt.Errorf("list repos: %w", err)
	}
	defer rows.Close()

	var repos []Repository
	for rows.Next() {
		var r Repository
		if err := rows.Scan(&r.ID, &r.OwnerID, &r.OwnerName, &r.Name, &r.Description, &r.IsPrivate, &r.DefaultBranch, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan repo: %w", err)
		}
		repos = append(repos, r)
	}
	return repos, rows.Err()
}

// ListPublic returns all public repositories.
func (s *RepoStore) ListPublic() ([]Repository, error) {
	rows, err := s.DB.Query(
		`SELECT r.id, r.owner_id, u.username, r.name, r.description, r.is_private, r.default_branch, r.created_at, r.updated_at
		 FROM repositories r JOIN users u ON r.owner_id = u.id
		 WHERE r.is_private = FALSE ORDER BY r.updated_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list public repos: %w", err)
	}
	defer rows.Close()

	var repos []Repository
	for rows.Next() {
		var r Repository
		if err := rows.Scan(&r.ID, &r.OwnerID, &r.OwnerName, &r.Name, &r.Description, &r.IsPrivate, &r.DefaultBranch, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan repo: %w", err)
		}
		repos = append(repos, r)
	}
	return repos, rows.Err()
}

// Delete removes a repository record.
func (s *RepoStore) Delete(id int64) error {
	_, err := s.DB.Exec(`DELETE FROM repositories WHERE id = ?`, id)
	return err
}
