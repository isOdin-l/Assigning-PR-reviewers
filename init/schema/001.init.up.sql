CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id TEXT PRIMARY KEY,
    team_name TEXT NOT NULL,
    user_name TEXT UNIQUE NOT NULL,
    is_active BOOLEAN NOT NULL
);

CREATE TYPE PR_STATUS AS ENUM ('OPEN', 'MERGED');
CREATE TABLE pull_requests(
    id TEXT PRIMARY KEY,
    author_id TEXT REFERENCES users(id) NOT NULL,
    name TEXT NOT NULL,
    status PR_STATUS DEFAULT 'OPEN' NOT NULL,
    merged_at TIMESTAMPTZ
    -- created_at
);


CREATE TABLE pr_reviewers(
    pr_id TEXT REFERENCES pull_requests(id) ON DELETE CASCADE,
    reviewer_id TEXT REFERENCES users(id),
    PRIMARY KEY(pr_id, reviewer_id)
);