CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS forums CASCADE;
DROP TABLE IF EXISTS threads CASCADE;
DROP TABLE IF EXISTS posts CASCADE;
DROP TABLE IF EXISTS votes CASCADE;

SET LOCAL synchronous_commit TO OFF;

CREATE TABLE users (
    id integer GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    email citext NOT NULL UNIQUE,
    nickname citext NOT NULL COLLATE "C" UNIQUE,
    fullname text DEFAULT NULL,
    about text DEFAULT NULL
);

CREATE TABLE forums (
    id integer GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    slug citext NOT NULL UNIQUE,
    title text NOT NULL,
    username citext NOT NULL REFERENCES users (nickname),
    threads integer DEFAULT 0,
    posts integer DEFAULT 0
);

CREATE TABLE threads (
    id integer GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    slug citext DEFAULT NULL,
    created timestamptz DEFAULT NOW(),
    title text NOT NULL,
    message text DEFAULT NULL,
    username citext NOT NULL REFERENCES users (nickname),
    forum citext NOT NULL REFERENCES forums (slug),
    votes integer DEFAULT 0
);

CREATE TABLE posts (
    id integer GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    created timestamptz DEFAULT NOW(),
    isedited boolean DEFAULT FALSE,
    message text DEFAULT NULL,
    username citext NOT NULL REFERENCES users (nickname),
    forum citext NOT NULL REFERENCES forums (slug),
    thread integer NOT NULL REFERENCES threads,
    parent integer NOT NULL,
    path integer[],
    root integer 
);

CREATE TABLE votes (
    id integer GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    thread_id integer NOT NULL REFERENCES threads,
    username citext NOT NULL REFERENCES users (nickname),
    value integer NOT NULL,
	CONSTRAINT unique_votes UNIQUE (thread_id, username)
);

-- ==========================

CREATE OR REPLACE FUNCTION create_thread() RETURNS TRIGGER AS
$thread_trigger$
	BEGIN
        UPDATE forums SET threads = threads + 1
	    WHERE slug = NEW.forum;
        RETURN NEW;
	END;
$thread_trigger$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION create_post() RETURNS TRIGGER AS
$post_trigger$
	BEGIN
		IF (NEW.parent = 0) THEN 
            NEW.path = ARRAY[NEW.id];
            NEW.root = NEW.id;
        ELSE 
            NEW.path = (SELECT posts.path || NEW.id FROM posts WHERE id = NEW.parent);
            NEW.root = NEW.path[1];
		END IF;
		
        UPDATE forums SET posts = posts + 1 
        WHERE slug = NEW.forum;
		RETURN NEW;
	END;
$post_trigger$
LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION create_vote() RETURNS TRIGGER AS
$create_vote_trigger$
    BEGIN
        UPDATE threads SET votes = votes + NEW.value
        WHERE id = NEW.thread_id;
        RETURN NEW;
    END;
$create_vote_trigger$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_vote() RETURNS TRIGGER AS
$update_vote_trigger$
	BEGIN
        IF (OLD.value = 1 AND NEW.value = -1) THEN
            UPDATE threads SET votes = votes - 2
            WHERE id = NEW.thread_id;
        END IF;
        IF (OLD.value = -1 AND NEW.value = 1) THEN
            UPDATE threads SET votes = votes + 2
            WHERE id = NEW.thread_id;
        END IF;
        RETURN NEW;
	END;
$update_vote_trigger$
LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS thread_trigger ON posts;
CREATE TRIGGER thread_trigger BEFORE INSERT ON threads FOR EACH ROW EXECUTE PROCEDURE create_thread();

DROP TRIGGER IF EXISTS post_trigger ON posts;
CREATE TRIGGER post_trigger BEFORE INSERT ON posts FOR EACH ROW EXECUTE PROCEDURE create_post();

DROP TRIGGER IF EXISTS create_vote_trigger ON votes;
CREATE TRIGGER create_vote_trigger BEFORE INSERT ON votes FOR EACH ROW EXECUTE PROCEDURE create_vote();

DROP TRIGGER IF EXISTS update_vote_triggerr ON votes;
CREATE TRIGGER update_vote_trigger BEFORE UPDATE ON votes FOR EACH ROW EXECUTE PROCEDURE update_vote();

-- CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
-- SELECT pg_stat_statements_reset();

-- SELECT * FROM pg_stat_statements Order by total_time, max_time;