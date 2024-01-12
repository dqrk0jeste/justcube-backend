CREATE TABLE users (
  id UUID PRIMARY KEY,
  username VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(60) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT(now())
);

CREATE TABLE posts (
  id UUID PRIMARY KEY,
  text_content VARCHAR(500) NOT NULL,
  image_count INT NOT NULL DEFAULT(0),
  user_id UUID NOT NULL REFERENCES users(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT(now())
);

CREATE TABLE comments (
  id UUID PRIMARY KEY,
  content VARCHAR(200) NOT NULL,
  user_id UUID NOT NULL REFERENCES users(id),
  post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT(now())
);

CREATE TABLE replies (
  id UUID PRIMARY KEY,
  content VARCHAR(200) NOT NULL,
  user_id UUID NOT NULL REFERENCES users(id),
  comment_id UUID NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT(now())
);

CREATE TABLE messages (
  id UUID PRIMARY KEY,
  content VARCHAR(500) NOT NULL,
  from_user_id UUID NOT NULL REFERENCES users(id),
  to_user_id UUID NOT NULL REFERENCES users(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT(now())
);

CREATE TABLE follows (
  user_id UUID NOT NULL REFERENCES users(id),
  followed_user_id UUID NOT NULL REFERENCES users(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT(now()),
  PRIMARY KEY(user_id, followed_user_id)
);

CREATE INDEX ON users (id);
CREATE INDEX ON posts (id);
CREATE INDEX ON comments (post_id);
CREATE INDEX ON messages (from_user_id);
CREATE INDEX ON messages (to_user_id);
CREATE INDEX ON messages (from_user_id, to_user_id);
CREATE INDEX ON follows (user_id);
CREATE INDEX ON follows (followed_user_id);
CREATE INDEX ON follows (user_id, followed_user_id);