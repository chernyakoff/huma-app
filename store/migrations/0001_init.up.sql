CREATE TABLE IF NOT EXISTS users (
  id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  email TEXT NOT NULL,
  password TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  role TEXT NOT NULL DEFAULT 'user' CHECK(role IN ('admin', 'user', 'editor')),
  UNIQUE (email)
);
