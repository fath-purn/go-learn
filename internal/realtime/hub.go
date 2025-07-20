package realtime

import (
	"encoding/json"
	"example/hello/internal/user"
	"log"
	"strconv"
	"time"
)

// Client adalah representasi dari satu pengguna yang terhubung melalui WebSocket.
type Client struct {
	// UserID adalah ID dari pengguna yang terotentikasi.
	UserID string

	// RoomID adalah ID dari room tempat client ini bergabung.
	RoomID string

	// Hub tempat client ini terdaftar.
	Hub *Hub

	// Koneksi WebSocket itu sendiri.
	Conn Conn

	// Channel buffer untuk pesan keluar.
	Send chan ChatMessage
}

// Message mendefinisikan struktur untuk pesan chat.
// Ini akan dikonversi ke/dari JSON saat dikirim melalui WebSocket.
type ChatMessage struct {
	// Type bisa digunakan untuk membedakan jenis pesan, misal: "chat_message", "user_joined", dll.
	Type string `json:"type"`
	// Content adalah isi pesan itu sendiri.
	Content string `json:"content"`
	// SenderID adalah ID pengguna yang mengirim pesan.
	SenderID string `json:"sender_id"`
	// name
	SenderName string `json:"sender_name"`
	// RoomID adalah ID dari room tujuan pesan ini.
	RoomID string `json:"room_id"`
	// Timestamp kapan pesan dibuat di server.
	Timestamp time.Time `json:"timestamp"`
}

// Conn adalah interface untuk koneksi WebSocket agar mudah di-mock saat testing.
type Conn interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	Close() error
}

// Hub mengelola semua client dan melakukan broadcast pesan.
type Hub struct {
	// Kumpulan room yang aktif. Kunci adalah RoomID, nilai adalah peta client di room itu.
	rooms map[string]map[*Client]bool

	// Pesan masuk dari client.
	Broadcast chan ChatMessage

	// Permintaan registrasi dari client.
	Register chan *Client

	// Permintaan unregistrasi dari client.
	Unregister chan *Client

	// Service untuk menyimpan pesan ke database.
	messageService Service

	// untuk data user
	userService user.Service
}

// NewHub membuat instance Hub baru.
func NewHub(messageService Service, userService user.Service) *Hub {
	return &Hub{
		Broadcast:      make(chan ChatMessage),
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		rooms:          make(map[string]map[*Client]bool),
		messageService: messageService,
		userService:    userService,
	}
}

// Run menjalankan Hub dalam sebuah goroutine.
// Ini adalah event loop utama untuk semua aktivitas real-time.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			// buat room
			if _, ok := h.rooms[client.RoomID]; !ok {
				h.rooms[client.RoomID] = make(map[*Client]bool)
			}

			// mendaftarkan client
			h.rooms[client.RoomID][client] = true

			// Jalankan goroutine untuk mengirim riwayat chat ke client yang baru terhubung.
			go h.sendChatHistory(client)
			log.Printf("Client %s terhubung ke room %s", client.UserID, client.RoomID)

		case client := <-h.Unregister:
			// hapus client dari room
			if room, ok := h.rooms[client.RoomID]; ok {
				if _, ok := room[client]; ok {
					delete(room, client)
					close(client.Send)
					log.Printf("Client %s terputus dari room %s. Sisa: %d", client.UserID, client.RoomID, len(room))
				}
				// jika room kosong, maka hapus room
				if len(room) == 0 {
					delete(h.rooms, client.RoomID)
					log.Printf("Room %s kosong. Dihapus.", client.RoomID)
				}
			}

		case chatMessage := <-h.Broadcast:
			// 1. Simpan pesan ke database sebelum broadcast
			senderID, err := strconv.ParseUint(chatMessage.SenderID, 10, 32)
			if err != nil {
				log.Printf("Gagal konversi SenderID '%s' ke uint: %v", chatMessage.SenderID, err)
				continue
			}

			dbMessage := Message{
				Content:  chatMessage.Content,
				RoomID:   chatMessage.RoomID,
				SenderID: uint(senderID),
			}
			_, err = h.messageService.SaveMessage(dbMessage)
			if err != nil {
				log.Printf("Gagal menyimpan pesan ke database: %v", err)
				// Kita tetap melanjutkan broadcast meskipun gagal menyimpan
			}

			// ambil nama user
			sender, err := h.userService.FindByID(int(senderID))
			if err == nil {
				chatMessage.SenderName = sender.Name
			}

			// Marshal pesan ke JSON untuk logging
			if msgJson, err := json.Marshal(chatMessage); err == nil {
				log.Printf("Broadcasting pesan ke room %s: %s", chatMessage.RoomID, msgJson)
			}

			// kirim pesan ke room sesuai
			if room, ok := h.rooms[chatMessage.RoomID]; ok {
				for client := range room {
					select {
					case client.Send <- chatMessage:
					default:
						// jika channel penuh maka hapus koneksi
						close(client.Send)
						delete(room, client)
					}
				}
			}
		}
	}
}

// sendChatHistory mengambil riwayat chat dari database dan mengirimkannya ke satu client.
func (h *Hub) sendChatHistory(client *Client) {
	history, err := h.messageService.GetMessageByRoom(client.RoomID)
	if err != nil {
		log.Printf("Gagal mengambil riwayat chat untuk room %s: %v", client.RoomID, err)
		return
	}

	log.Printf("Mengirim %d pesan riwayat dari room %s ke client %s", len(history), client.RoomID, client.UserID)

	for _, msg := range history {
		senderName := ""
		sender, err := h.userService.FindByID(int(msg.SenderID))
		if err == nil {
			senderName = sender.Name
		}

		chatMessage := ChatMessage{
			Type:       "history_message", // Tipe pesan khusus untuk riwayat
			Content:    msg.Content,
			RoomID:     msg.RoomID,
			SenderID:   strconv.FormatUint(uint64(msg.SenderID), 10),
			SenderName: senderName,
			Timestamp:  msg.CreatedAt,
		}

		// Kirim pesan ke channel Send milik client.
		// Ini akan diambil oleh writePump client tersebut.
		select {
		case client.Send <- chatMessage:
		default:
			log.Printf("Channel send untuk client %s penuh, membatalkan pengiriman riwayat.", client.UserID)
			return
		}
	}
}
