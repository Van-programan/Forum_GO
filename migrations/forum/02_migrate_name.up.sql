CREATE TABLE IF NOT EXISTS users
(
	id SERIAL PRIMARY KEY,
	username VARCHAR(50) NOT NULL UNIQUE,
	role TEXT NOT NULL DEFAULT 'user',
	password_hash TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS topics (
    id SERIAL PRIMARY KEY,
    category_id INT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    author_id INT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    topic_id INT NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    author_id INT REFERENCES users(id) ON DELETE SET NULL,
    content TEXT NOT NULL,
    reply_to INT REFERENCES posts(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    user_id integer,
    username text NOT NULL,
    content text NOT NULL,
    created_at timestamp with time zone
);

CREATE INDEX IF NOT EXISTS idx_topics_category_id ON public.topics(category_id);

CREATE INDEX IF NOT EXISTS idx_topics_author_id ON public.topics(author_id);

CREATE INDEX IF NOT EXISTS idx_posts_topic_id ON public.posts(topic_id);

CREATE INDEX IF NOT EXISTS idx_posts_author_id ON public.posts(author_id);

CREATE INDEX IF NOT EXISTS idx_posts_reply_to ON public.posts(reply_to);

