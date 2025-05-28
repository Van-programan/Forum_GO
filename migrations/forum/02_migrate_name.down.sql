DROP TABLE IF EXISTS categories;

DROP TABLE IF EXISTS topics;

DROP TABLE IF EXISTS posts;

DROP TABLE IF EXISTS messages;

DROP INDEX IF EXISTS idx_topics_category_id;

DROP INDEX IF EXISTS idx_topics_author_id;

DROP INDEX IF EXISTS idx_posts_topic_id;

DROP INDEX IF EXISTS idx_posts_author_id;

DROP INDEX IF EXISTS idx_posts_reply_to;