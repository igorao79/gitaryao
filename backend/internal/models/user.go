package models

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	AvatarURL string    `json:"avatar_url"`
	GithubID  *int64    `json:"github_id,omitempty"`
	GoogleID  *string   `json:"google_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserStore struct {
	DB *sql.DB
}

// FindByGithubID finds a user by their GitHub ID.
func (s *UserStore) FindByGithubID(githubID int64) (*User, error) {
	u := &User{}
	err := s.DB.QueryRow(
		`SELECT id, username, email, avatar_url, github_id, google_id, created_at, updated_at
		 FROM users WHERE github_id = ?`, githubID,
	).Scan(&u.ID, &u.Username, &u.Email, &u.AvatarURL, &u.GithubID, &u.GoogleID, &u.CreatedAt, &u.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find by github id: %w", err)
	}
	return u, nil
}

// FindByGoogleID finds a user by their Google ID.
func (s *UserStore) FindByGoogleID(googleID string) (*User, error) {
	u := &User{}
	err := s.DB.QueryRow(
		`SELECT id, username, email, avatar_url, github_id, google_id, created_at, updated_at
		 FROM users WHERE google_id = ?`, googleID,
	).Scan(&u.ID, &u.Username, &u.Email, &u.AvatarURL, &u.GithubID, &u.GoogleID, &u.CreatedAt, &u.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find by google id: %w", err)
	}
	return u, nil
}

// FindByID finds a user by their database ID.
func (s *UserStore) FindByID(id int64) (*User, error) {
	u := &User{}
	err := s.DB.QueryRow(
		`SELECT id, username, email, avatar_url, github_id, google_id, created_at, updated_at
		 FROM users WHERE id = ?`, id,
	).Scan(&u.ID, &u.Username, &u.Email, &u.AvatarURL, &u.GithubID, &u.GoogleID, &u.CreatedAt, &u.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find by id: %w", err)
	}
	return u, nil
}

// FindByUsername finds a user by their username.
func (s *UserStore) FindByUsername(username string) (*User, error) {
	u := &User{}
	err := s.DB.QueryRow(
		`SELECT id, username, email, avatar_url, github_id, google_id, created_at, updated_at
		 FROM users WHERE username = ?`, username,
	).Scan(&u.ID, &u.Username, &u.Email, &u.AvatarURL, &u.GithubID, &u.GoogleID, &u.CreatedAt, &u.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find by username: %w", err)
	}
	return u, nil
}

// UpsertGithubUser creates or updates a user from GitHub OAuth data.
func (s *UserStore) UpsertGithubUser(githubID int64, username, email, avatarURL string) (*User, error) {
	existing, err := s.FindByGithubID(githubID)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		// Update existing user
		_, err := s.DB.Exec(
			`UPDATE users SET username = ?, email = ?, avatar_url = ?, updated_at = CURRENT_TIMESTAMP
			 WHERE github_id = ?`,
			username, email, avatarURL, githubID,
		)
		if err != nil {
			return nil, fmt.Errorf("update github user: %w", err)
		}
		return s.FindByGithubID(githubID)
	}

	// Create new user
	result, err := s.DB.Exec(
		`INSERT INTO users (username, email, avatar_url, github_id) VALUES (?, ?, ?, ?)`,
		username, email, avatarURL, githubID,
	)
	if err != nil {
		return nil, fmt.Errorf("insert github user: %w", err)
	}

	id, _ := result.LastInsertId()
	return s.FindByID(id)
}

// UpsertGoogleUser creates or updates a user from Google OAuth data.
func (s *UserStore) UpsertGoogleUser(googleID, username, email, avatarURL string) (*User, error) {
	existing, err := s.FindByGoogleID(googleID)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		_, err := s.DB.Exec(
			`UPDATE users SET email = ?, avatar_url = ?, updated_at = CURRENT_TIMESTAMP
			 WHERE google_id = ?`,
			email, avatarURL, googleID,
		)
		if err != nil {
			return nil, fmt.Errorf("update google user: %w", err)
		}
		return s.FindByGoogleID(googleID)
	}

	result, err := s.DB.Exec(
		`INSERT INTO users (username, email, avatar_url, google_id) VALUES (?, ?, ?, ?)`,
		username, email, avatarURL, googleID,
	)
	if err != nil {
		return nil, fmt.Errorf("insert google user: %w", err)
	}

	id, _ := result.LastInsertId()
	return s.FindByID(id)
}
