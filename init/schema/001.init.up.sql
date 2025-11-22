CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE teams(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(256) NOT NULL,
    team_name TEXT REFERENCES teams(name) NOT NULL,
    is_active BOOLEAN NOT NULL
);


CREATE TYPE PR_STATUS AS ENUM ('OPEN', 'MERGED');
CREATE TABLE pull_requests(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    author_id UUID REFERENCES users(id) NOT NULL,
    name TEXT NOT NULL,
    status PR_STATUS DEFAULT 'OPEN' NOT NULL
);


CREATE TABLE pr_reviewers(
    pr_id UUID REFERENCES pull_requests(id) ON DELETE CASCADE,
    reviewer_id UUID REFERENCES users(id),
    PRIMARY KEY(pr_id, reviewer_id)
);

CREATE OR REPLACE FUNCTION reviewers_limit() RETURNS TRIGGER AS $$ DECLARE cnt INT;
BEGIN
    SELECT COUNT(*) INTO cnt FROM pr_reviewers WHERE pr_id = NEW.pr_id;

    IF cnt >= 2 THEN
        RAISE EXCEPTION 'PR already has % 2 reviewers', NEW.pr_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER pr_reviewers_limit BEFORE INSERT ON pr_reviewers FOR EACH ROW EXECUTE FUNCTION reviewers_limit();