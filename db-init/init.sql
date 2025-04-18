\getenv dbname POSTGRES_DBNAME
CREATE DATABASE :dbname WITH ENCODING='UTF8';

\connect :dbname

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    handle VARCHAR(255) NOT NULL UNIQUE,
    is_admin BOOLEAN NOT NULL,
    auth_bits VARCHAR(255) NOT NULL
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    author_id INTEGER NOT NULL REFERENCES users(id),
    upload_timestamp TIMESTAMP NOT NULL,
    modified_timestamp TIMESTAMP NOT NULL,
    release_timestamp TIMESTAMP NOT NULL,
    archived BOOLEAN NOT NULL,
    source_ref VARCHAR(255) NOT NULL
);

CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE tags_to_posts (
    tag_id INTEGER NOT NULL REFERENCES tags(id),
    post_id INTEGER NOT NULL REFERENCES posts(id)
);

CREATE TABLE action_log (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    ts TIMESTAMP NOT NULL,
    action VARCHAR(255) NOT NULL
);
