package gatewayserver

import (
	"encoding/json"
	"log"
	"quiz/internal/constant"
	"quiz/models"
	"quiz/pkg/mkafka"
	"sync"

	"github.com/gorilla/websocket"
)

// a hub use to manage connection and room of all user
type Hub struct {
	mu    sync.RWMutex
	Users map[string]*websocket.Conn
	Rooms map[string]map[string]struct{} // key room is quizId, and value is all user in that room
}

func (h *Hub) run() {
	consumer, err := mkafka.ConsumeMessages(constant.KAFKATOPIC_LEADERBOARDCHANGE)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-consumer.Errors():
			log.Println(err)
		case msg := <-consumer.Messages():
			var data models.LeaderboardMemberChange
			if err := json.Unmarshal(msg.Value, &data); err != nil {
				log.Println(err)
				continue
			}

			// broadcast to all user in room
			go h.broadcastMemberChange(data)
		}
	}
}

func (h *Hub) broadcastMemberChange(data models.LeaderboardMemberChange) {
	userIDs, exist := h.Rooms[data.QuizID]
	if !exist { // no room in this server
		return
	}

	msg := WebsocketMessage[models.LeaderboardMemberChange]{
		Type:    0,
		Message: data,
	}
	msgByte, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()
	for userID := range userIDs {
		conn, exist := h.Users[userID]
		if exist {
			go conn.WriteMessage(websocket.BinaryMessage, msgByte)
		}
	}
}

func (h *Hub) joinRoom(userID, roomID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// set connection
	h.Users[userID] = conn

	// insert to room
	room, exist := h.Rooms[roomID]
	if !exist {
		room = make(map[string]struct{})
		h.Rooms[roomID] = room
	}
	room[userID] = struct{}{}
}

func (h *Hub) leaveRoom(userID, roomID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// delete connection key
	delete(h.Users, userID)

	// leave to room
	room, exist := h.Rooms[roomID]
	if !exist || room == nil {
		return
	}
	delete(room, userID)

	// delete room if empty
	if len(room) == 0 {
		delete(h.Rooms, roomID)
	}
}
