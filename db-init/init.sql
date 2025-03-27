\getenv dbname POSTGRES_DBNAME
CREATE DATABASE :dbname WITH ENCODING='UTF8';

\connect :dbname

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    handle VARCHAR(255),
    is_admin BOOLEAN,
    auth_bits INTEGER NOT NULL
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    author_id INTEGER REFERENCES users(id),
    upload_timestamp TIMESTAMP,
    modified_timestamp TIMESTAMP,
    release_timestamp TIMESTAMP,
    archived BOOLEAN,
    source_ref VARCHAR(255)
);

CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255)
);

CREATE TABLE tags_to_posts (
    tag_id INTEGER REFERENCES tags(id),
    post_id INTEGER REFERENCES posts(id)
);

CREATE TABLE action_log (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    ts TIMESTAMP,
    action VARCHAR(255)
);
