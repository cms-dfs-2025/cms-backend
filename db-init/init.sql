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

    draft BOOLEAN NOT NULL,
    archived BOOLEAN NOT NULL,

    title VARCHAR(255) NOT NULL,
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
