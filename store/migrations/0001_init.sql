-- +goose Up
CREATE TABLE IF NOT EXISTS users (
  id TEXT NOT NULL PRIMARY KEY,
  email TEXT NOT NULL,
  password TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  role TEXT NOT NULL DEFAULT 'user' CHECK(role IN ('admin', 'user', 'editor')),
  verified INTEGER NOT NULL DEFAULT 0 CHECK(verified IN (0,1)),
  UNIQUE (email)
);

-- +goose Down
DROP TABLE IF EXISTS users;