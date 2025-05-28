package response

import "github.com/Van-programan/Forum_GO/internal/entity"

type ErrorResponseForum struct {
	Error string `json:"error" example:"error message"`
}

type SuccessMessageResponse struct {
	Message string `json:"message" example:"operation was successful"`
}

type IDResponse struct {
	ID int64 `json:"id" example:"123"`
}

type CategoryResponse struct {
	Category entity.Category `json:"category"`
}

type CategoriesResponse struct {
	Categories []entity.Category `json:"categories"`
}

type TopicResponse struct {
	Topic entity.Topic `json:"topic"`
}

type TopicsResponse struct {
	Topics []entity.Topic `json:"topics"`
}

type PostsResponse struct {
	Posts []entity.Post `json:"posts"`
}
