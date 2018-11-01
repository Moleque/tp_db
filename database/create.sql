CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS forums;
DROP TABLE IF EXISTS threads;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS votes;


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
    username citext NOT NULL,
    threads integer DEFAULT 0,
    posts integer DEFAULT 0
);

CREATE TABLE threads (
    id integer GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    slug citext DEFAULT NULL,
    created timestamptz DEFAULT NOW(),
    title text NOT NULL,
    message text DEFAULT NULL,
    username citext NOT NULL,
    forum citext NOT NULL,
    votes integer DEFAULT 0
);

CREATE TABLE posts (
    id integer GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    created timestamptz DEFAULT NOW(),
    isedited boolean DEFAULT FALSE,
    message text DEFAULT NULL,
    username citext NOT NULL,
    forum citext NOT NULL,
    thread integer NOT NULL,
    parent integer DEFAULT NULL,
    path integer[],
    root integer 
);

CREATE TABLE votes (
    id integer GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    thread_id integer NOT NULL,
    username citext NOT NULL,
    value integer NOT NULL
);


CREATE OR REPLACE FUNCTION create_path() RETURNS TRIGGER AS
$path_trigger$
	BEGIN
		IF (NEW.parent = 0)
			THEN 
				NEW.path = ARRAY[NEW.id];
				NEW.root = NEW.id;
			ELSE 
				NEW.path = (SELECT posts.path || NEW.id FROM posts WHERE id = NEW.parent);
				NEW.root = NEW.path[1];
		END IF;
		UPDATE forums SET posts = posts + 1 WHERE slug = NEW.forum;
		RETURN NEW;
	END;
$path_trigger$
LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS path_trigger ON posts;
CREATE TRIGGER path_trigger BEFORE INSERT ON posts FOR EACH ROW EXECUTE PROCEDURE create_path();