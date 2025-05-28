package entity

import "time"

type ChatMessage struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Category struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Post struct {
	ID        int64     `json:"id"`
	TopicID   int64     `json:"topic_id"`
	AuthorID  *int64    `json:"author_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	ReplyTo   *int64    `json:"reply_to"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Topic struct {
	ID         int64     `json:"id"`
	CategoryID int64     `json:"category_id"`
	Title      string    `json:"title"`
	AuthorID   *int64    `json:"author_id"`
	Username   string    `json:"username"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type WsMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
}

type IncomingWsMessage struct {
	Content string `json:"content"`
}
