package transport

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/Van-programan/Forum_GO/internal/transport/ws"
	"github.com/Van-programan/Forum_GO/internal/usecase"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type WSHandler struct {
	forumUC  usecase.Forum
	upgrader websocket.Upgrader
}

func NewWSHandler(forumUC usecase.Forum) *WSHandler {
	return &WSHandler{
		forumUC: forumUC,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (h *WSHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/ws/{topicID}", h.HandleWebSocket)
}

func (h *WSHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {

	token := r.URL.Query().Get("token")
	userID, err := h.forumUC.ValidateToken(token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	topicID := vars["topicID"]
	if !h.forumUC.CanAccessTopic(userID, topicID) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	client := &ws.Client{
		ID:     generateClientID(),
		UserID: userID,
		Conn:   conn,
		Send:   make(chan ws.Message, 256),
	}

	h.forumUC.RegisterWSClient(topicID, client)

	go client.WritePump()
}

func generateClientID() string {
	return strconv.FormatInt(int64(rand.Intn(10000)), 10)
}
