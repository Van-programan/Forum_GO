package ws

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type HubInterface interface {
	RegisterClient(client *Client)
	UnregisterClient(client *Client)
	BroadcastToTopic(topicID string, msg Message)
}
