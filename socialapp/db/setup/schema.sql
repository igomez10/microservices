CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE SEQUENCE IF NOT EXISTS users_id_seq;
CREATE TABLE IF NOT EXISTS users (
	id BIGINT NOT NULL PRIMARY KEY DEFAULT nextval('users_id_seq'::regclass),
	username VARCHAR(100) NOT NULL,
	hashed_password VARCHAR(256) NOT NULL DEFAULT '',
	hashed_password_expires_at TIMESTAMP NOT NULL DEFAULT '2000-01-01 00:00:00',
	salt VARCHAR(100) NOT NULL DEFAULT '',
	first_name VARCHAR(100) NOT NULL,
	last_name VARCHAR(100) NOT NULL,
	email VARCHAR(100) NOT NULL,
	email_token VARCHAR(100) NOT NULL DEFAULT '',
	email_verified_at TIMESTAMP NULL DEFAULT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP NULL
);
CREATE INDEX IF NOT EXISTS users_username_idx ON users (username);
CREATE INDEX IF NOT EXISTS users_email_idx ON users (email);
CREATE INDEX IF NOT EXISTS users_deleted_at_idx ON users (deleted_at);
CREATE INDEX IF NOT EXISTS users_created_at_idx ON users (created_at);
CREATE INDEX IF NOT EXISTS users_username_trgm_idx ON users USING GIN (username gin_trgm_ops);
CREATE INDEX IF NOT EXISTS users_first_name_trgm_idx ON users USING GIN (first_name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS users_last_name_trgm_idx ON users USING GIN (last_name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS users_email_trgm_idx ON users USING GIN (email gin_trgm_ops);

CREATE SEQUENCE IF NOT EXISTS comments_id_seq;
CREATE TABLE IF NOT EXISTS comments (
    id BIGINT NOT NULL PRIMARY KEY DEFAULT nextval('comments_id_seq'::regclass),
    content VARCHAR(8192) NOT NULL,
    like_count BIGINT NOT NULL DEFAULT 0,
    user_id BIGINT NOT NULL,
    CONSTRAINT fk_user_id
      FOREIGN KEY(user_id)
	  REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE INDEX IF NOT EXISTS comments_user_id_idx ON comments (user_id);
CREATE INDEX IF NOT EXISTS comments_deleted_at_idx ON comments (deleted_at);
CREATE INDEX IF NOT EXISTS comments_created_at_idx ON comments (created_at);
CREATE INDEX IF NOT EXISTS comments_user_id_created_at_idx ON comments (user_id, created_at DESC);

CREATE SEQUENCE IF NOT EXISTS likes_id_seq;
CREATE TABLE IF NOT EXISTS likes (
    id BIGINT NOT NULL PRIMARY KEY DEFAULT nextval('likes_id_seq'::regclass),
    user_id BIGINT NOT NULL,
    comment_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT likes_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT likes_comment_id_fkey FOREIGN KEY (comment_id) REFERENCES comments(id)
);
CREATE INDEX IF NOT EXISTS likes_user_id_idx ON likes (user_id);
CREATE INDEX IF NOT EXISTS likes_comment_id_idx ON likes (comment_id);
CREATE INDEX IF NOT EXISTS likes_deleted_at_idx ON likes (deleted_at);
CREATE INDEX IF NOT EXISTS likes_comment_id_deleted_at_idx ON likes (comment_id, deleted_at);
CREATE UNIQUE INDEX IF NOT EXISTS likes_user_id_comment_id_active_uidx
ON likes (user_id, comment_id)
WHERE deleted_at IS NULL;

CREATE SEQUENCE IF NOT EXISTS followers_id_seq;
CREATE TABLE IF NOT EXISTS followers (
    id BIGINT NOT NULL DEFAULT nextval('followers_id_seq'::regclass) PRIMARY KEY,
    follower_id BIGINT NOT NULL,
    followed_id BIGINT NOT NULL,
    CONSTRAINT followers_follower_id_fkey FOREIGN KEY ("follower_id") REFERENCES "public"."users"("id"),
    CONSTRAINT followers_followed_id_fkey FOREIGN KEY ("followed_id") REFERENCES "public"."users"("id"),
    UNIQUE (follower_id, followed_id),
    UNIQUE (follower_id, followed_id)
);
CREATE INDEX IF NOT EXISTS followers_follower_id_idx ON followers (follower_id);
CREATE INDEX IF NOT EXISTS followers_followed_id_idx ON followers (followed_id);
CREATE INDEX IF NOT EXISTS followers_follower_id_followed_id_idx ON followers (follower_id, followed_id);
CREATE INDEX IF NOT EXISTS followers_followed_id_follower_id_idx ON followers (followed_id, follower_id);

CREATE SEQUENCE IF NOT EXISTS credentials_id_seq;
CREATE TABLE IF NOT EXISTS credentials ( -- long term api keys
    id BIGINT NOT NULL DEFAULT nextval('credentials_id_seq'::regclass) PRIMARY KEY,
    user_id BIGINT NOT NULL,
    public_key VARCHAR(512) NOT NULL,
    description VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT credentials_user_id_fkey FOREIGN KEY (user_id) REFERENCES "public"."users"("id"),
    UNIQUE (user_id, public_key)
);
CREATE INDEX IF NOT EXISTS credentials_user_id_idx ON credentials (user_id);
CREATE INDEX IF NOT EXISTS credentials_public_key_idx ON credentials (public_key);
CREATE INDEX IF NOT EXISTS credentials_deleted_at_idx ON credentials (deleted_at);

CREATE SEQUENCE IF NOT EXISTS roles_id_seq;
CREATE TABLE IF NOT EXISTS roles (
    id BIGINT NOT NULL DEFAULT nextval('roles_id_seq'::regclass) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
CREATE INDEX IF NOT EXISTS roles_name_idx ON roles (name);

CREATE SEQUENCE IF NOT EXISTS scopes_id_seq;
CREATE TABLE IF NOT EXISTS scopes (
    id BIGINT NOT NULL DEFAULT nextval('scopes_id_seq'::regclass) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS scopes_name_idx ON scopes (name);

CREATE SEQUENCE IF NOT EXISTS roles_to_scopes_id_seq;
CREATE TABLE IF NOT EXISTS roles_to_scopes (
    id BIGINT NOT NULL DEFAULT nextval('roles_to_scopes_id_seq'::regclass) PRIMARY KEY,
    role_id BIGINT NOT NULL,
    scope_id BIGINT NOT NULL,
    CONSTRAINT roles_to_scopes_role_id_fkey FOREIGN KEY (role_id) REFERENCES "public"."roles"("id"),
    CONSTRAINT roles_to_scopes_scope_id_fkey FOREIGN KEY (scope_id) REFERENCES "public"."scopes"("id"),
    UNIQUE (role_id, scope_id)
);
CREATE INDEX IF NOT EXISTS roles_to_scopes_role_id_idx ON roles_to_scopes (role_id);
CREATE INDEX IF NOT EXISTS roles_to_scopes_scope_id_idx ON roles_to_scopes (scope_id);
CREATE INDEX IF NOT EXISTS roles_to_scopes_role_id_scope_id_idx ON roles_to_scopes (role_id, scope_id);
CREATE INDEX IF NOT EXISTS roles_to_scopes_scope_id_role_id_idx ON roles_to_scopes (scope_id, role_id);

CREATE SEQUENCE IF NOT EXISTS users_to_roles_id_seq;
CREATE TABLE IF NOT EXISTS users_to_roles (
    id BIGINT NOT NULL DEFAULT nextval('users_to_roles_id_seq'::regclass) PRIMARY KEY,
    role_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    UNIQUE (role_id, user_id),
    CONSTRAINT users_to_roles_role_id_fkey FOREIGN KEY (role_id) REFERENCES "public"."roles"("id"),
    CONSTRAINT users_to_roles_user_id_fkey FOREIGN KEY (user_id) REFERENCES "public"."users"("id")
);
CREATE INDEX IF NOT EXISTS users_to_roles_role_id_idx ON users_to_roles (role_id);
CREATE INDEX IF NOT EXISTS users_to_roles_user_id_idx ON users_to_roles (user_id);
CREATE INDEX IF NOT EXISTS users_to_roles_role_id_user_id_idx ON users_to_roles (role_id, user_id);
CREATE INDEX IF NOT EXISTS users_to_roles_user_id_role_id_idx ON users_to_roles (user_id, role_id);

-- Events table for event sourcing
CREATE SEQUENCE IF NOT EXISTS events_id_seq;
CREATE TABLE IF NOT EXISTS events (
    id BIGINT NOT NULL DEFAULT nextval('events_id_seq'::regclass) PRIMARY KEY,
    aggregate_id BIGINT NOT NULL,
    aggregate_type VARCHAR(255) NOT NULL,
    version BIGINT NOT NULL,
    event_type VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
CREATE INDEX IF NOT EXISTS events_aggregate_id_idx ON events (aggregate_id);
CREATE INDEX IF NOT EXISTS events_aggregate_type_idx ON events (aggregate_type);
-- Each aggregate's version is unique: enforces a strict, gap-checkable ordering
-- of events per entity and provides optimistic-concurrency protection.
CREATE UNIQUE INDEX IF NOT EXISTS events_aggregate_version_idx ON events (aggregate_id, aggregate_type, version);
CREATE INDEX IF NOT EXISTS events_version_idx ON events (version);
CREATE INDEX IF NOT EXISTS events_event_type_idx ON events (event_type);
CREATE INDEX IF NOT EXISTS events_created_at_idx ON events (created_at);


-- SHORTLY
CREATE SEQUENCE IF NOT EXISTS urls_id_seq;
CREATE TABLE IF NOT EXISTS urls (
	id BIGINT NOT NULL PRIMARY KEY DEFAULT nextval('urls_id_seq'::regclass),
	url VARCHAR(10000) NOT NULL,
	alias VARCHAR(10) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP NULL
);
CREATE INDEX IF NOT EXISTS urls_alias_idx ON urls (alias);



-- Change Data Capture (Apache Flink CDC): the Postgres source needs complete
-- before-images for UPDATE/DELETE events, otherwise non-key columns arrive NULL
-- and Flink CDC's schema inference fails. REPLICA IDENTITY FULL includes every
-- column in the WAL for the captured tables. Keep in sync with flink/pipeline.yaml.
ALTER TABLE events REPLICA IDENTITY FULL;
