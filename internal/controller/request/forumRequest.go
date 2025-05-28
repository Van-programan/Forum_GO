package request

type UpdateRequestCategory struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateRequestPost struct {
	Content string `json:"content"`
}

type UpdateRequestTopic struct {
	Title string `json:"title"`
}
