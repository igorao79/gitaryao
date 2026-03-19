package database

const schema = `
CREATE TABLE IF NOT EXISTS users (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    username    TEXT    NOT NULL UNIQUE,
    email       TEXT    NOT NULL UNIQUE,
    avatar_url  TEXT    DEFAULT '',
    github_id   INTEGER UNIQUE,
    google_id   TEXT    UNIQUE,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS repositories (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    owner_id        INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name            TEXT    NOT NULL,
    description     TEXT    DEFAULT '',
    is_private      BOOLEAN DEFAULT FALSE,
    default_branch  TEXT    DEFAULT 'master',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(owner_id, name)
);

CREATE TABLE IF NOT EXISTS ssh_keys (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title       TEXT    NOT NULL,
    public_key  TEXT    NOT NULL UNIQUE,
    fingerprint TEXT    NOT NULL,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_repositories_owner ON repositories(owner_id);
CREATE INDEX IF NOT EXISTS idx_ssh_keys_user ON ssh_keys(user_id);
`
