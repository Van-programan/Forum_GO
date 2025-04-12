package entity

import "time"

type Message struct {
	ID        int64     `json:"id" db:"id"`
	TopicID   int64     `json:"topic_id" db:"topic_id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Topic struct {
	ID        int64     `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	AuthorID  int64     `json:"author_id" db:"author_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
