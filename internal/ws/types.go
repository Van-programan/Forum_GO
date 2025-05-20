package ws

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type ClientRegister struct {
	UserID  string
	TopicID string
	Client  *Client
}

type HubEvent struct {
	Type string
	Data interface{}
}
