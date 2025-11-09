CREATE TABLE IF NOT EXISTS posts (
	id BIGINT PRIMARY KEY generated always AS identity,
	author_id uuid NOT NULL,
	title text NOT NULL,
	content text NOT NULL,
	created_at timestamptz NOT NULL DEFAULT NOW(),
	rating int NOT NULL DEFAULT 0,
	comments_restricted boolean NOT NULL DEFAULT false,
);

CREATE TABLE IF NOT EXISTS comments (
	id BIGINT PRIMARY KEY generated always AS identity,
	post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
	author_id uuid NOT NULL,
	"text" text NOT NULL,
	created_at timestamptz NOT NULL DEFAULT NOW(),
	rating int NOT NULL DEFAULT 0,
	deleted boolean NOT NULL DEFAULT false,
	parent_id bigint REFERENCES comments(id) ON DELETE CASCADE,
);

CREATE TABLE IF NOT EXISTS post_votes (
	voter_id uuid NOT NULL,
	post_id bigint NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
	value smallint NOT NULL,
	PRIMARY KEY (voter_id, post_id),
	CHECK (value IN (1, -1)),
);

CREATE TABLE IF NOT EXISTS comment_votes (
	voter_id uuid NOT NULL,
	comment_id bigint NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
	value smallint NOT NULL,
	PRIMARY KEY (voter_id, comment_id),
	CHECK (value IN (1, -1)),
);

CREATE INDEX posts_rating_id_idx ON posts (rating DESC, id ASC);

CREATE INDEX posts_created_at_desc_id_idx ON posts (created_at DESC, id ASC);

CREATE INDEX posts_created_at_asc_id_idx ON posts (created_at ASC, id ASC);

CREATE INDEX comments_post_id_rating_id_idx ON comments (post_id, rating DESC, id ASC);

CREATE INDEX comments_post_id_created_at_desc_id_idx ON comments (post_id, created_at DESC, id ASC);

CREATE INDEX comments_post_id_created_at_asc_id_idx ON comments (post_id DESC, created_at ASC, id ASC);

CREATE INDEX comments_post_id_idx ON comments (post_id);

CREATE INDEX comments_parent_id_idx ON comments (parent_id);