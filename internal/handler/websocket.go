package handler

import (
	"encoding/json"
	"example/hello/internal/realtime"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	hub *realtime.Hub
}

func NewWebSocketHandler(hub *realtime.Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Di produksi, validasi origin ini dengan benar.
		// Contoh: return r.Header.Get("Origin") == "http://example.com"
		return true
	},
}

// ServeWs menangani permintaan upgrade ke WebSocket.
func (h *WebSocketHandler) ServeWs(c *gin.Context) {
	// Middleware sudah berjalan dan memvalidasi token.
	// Kita hanya perlu mengambil UserID dari context.
	userIDVal, exists := c.Get("userID")
	if !exists {
		// Ini seharusnya tidak terjadi jika middleware dikonfigurasi dengan benar.
		log.Println("Error: userID tidak ditemukan di context Gin")
		return // Middleware seharusnya sudah menghentikan request ini.
	}

	userID, ok := userIDVal.(string)
	if !ok {
		log.Println("Error: userID di context bukan string")
		return
	}

	roomID := c.DefaultQuery("room", "general")

	// validasi untuk private room
	if strings.HasPrefix(roomID, "private-") {
		// memastikan user terontetikasi
		isValidParticipant, err := isUserParticipantOfPrivateRoom(userID, roomID)
		if err != nil || !isValidParticipant {
			log.Printf("Akses ditolak: User %s tidak valid untuk room private %s", userID, roomID)
			return
		}
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Gagal upgrade ke websocket:", err)
		return
	}

	// Gunakan UserID yang didapat dari token JWT di middleware.
	client := &realtime.Client{UserID: userID, RoomID: roomID, Hub: h.hub, Conn: conn, Send: make(chan realtime.ChatMessage, 256)}
	h.hub.Register <- client

	go h.writePump(client)
	go h.readPump(client)
}

func isUserParticipantOfPrivateRoom(currentUserID, roomID string) (bool, error) {
	parts := strings.Split(roomID, "-")
	if len(parts) != 3 {
		return false, fmt.Errorf("format private roomID tidak valid")
	}

	user1ID := parts[1]
	user2ID := parts[2]

	if currentUserID == user1ID || currentUserID == user2ID {
		return true, nil
	}

	return false, nil
}

func (h *WebSocketHandler) readPump(c *realtime.Client) {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	// Definisikan struktur untuk pesan masuk dari klien.
	// Klien hanya perlu mengirim kontennya.
	type incomingMessage struct {
		Content string `json:"content"`
	}

	for {
		// Baca pesan mentah (raw message) dari koneksi WebSocket
		_, rawMessage, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error saat membaca pesan: %v", err)
			}
			break
		}

		// Unmarshal pesan JSON yang masuk ke dalam struct
		var msg incomingMessage
		if err := json.Unmarshal(rawMessage, &msg); err != nil {
			log.Printf("Error unmarshaling pesan dari client %s: %v", c.UserID, err)
			continue // Lanjutkan ke pesan berikutnya jika ada format yang salah
		}

		// Buat pesan yang lengkap dengan data dari server (SenderID, Timestamp)
		fullMessage := realtime.ChatMessage{
			Type:      "chat_message",
			Content:   msg.Content,
			RoomID:    c.RoomID,
			SenderID:  c.UserID,
			Timestamp: time.Now(),
		}

		// Kirim pesan yang sudah terstruktur ke Hub untuk di-broadcast
		c.Hub.Broadcast <- fullMessage
	}
}

func (h *WebSocketHandler) writePump(c *realtime.Client) {
	defer c.Conn.Close()
	for message := range c.Send {
		// Marshal struct pesan menjadi JSON
		jsonMessage, err := json.Marshal(message)
		if err != nil {
			log.Printf("Error marshaling pesan untuk client %s: %v", c.UserID, err)
			continue
		}

		// Tulis pesan JSON ke koneksi WebSocket
		if err := c.Conn.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
			log.Printf("Error saat menulis pesan: %v", err)
			return
		}
	}
}
