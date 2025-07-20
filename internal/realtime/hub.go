package realtime

import (
	"encoding/json"
	"log"
	"strconv"
	"time"
)

// Client adalah representasi dari satu pengguna yang terhubung melalui WebSocket.
type Client struct {
	// UserID adalah ID dari pengguna yang terotentikasi.
	UserID string

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
	// Kumpulan client yang terdaftar.
	clients map[*Client]bool

	// Pesan masuk dari client.
	Broadcast chan ChatMessage

	// Permintaan registrasi dari client.
	Register chan *Client

	// Permintaan unregistrasi dari client.
	Unregister chan *Client

	// Service untuk menyimpan pesan ke database.
	messageService Service
}

// NewHub membuat instance Hub baru.
func NewHub(messageService Service) *Hub {
	return &Hub{
		Broadcast:      make(chan ChatMessage),
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		clients:        make(map[*Client]bool),
		messageService: messageService,
	}
}

// Run menjalankan Hub dalam sebuah goroutine.
// Ini adalah event loop utama untuk semua aktivitas real-time.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
			// Jalankan goroutine untuk mengirim riwayat chat ke client yang baru terhubung.
			go h.sendChatHistory(client)
			log.Println("Client baru terhubung. Total:", len(h.clients))

		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				log.Println("Client terputus. Sisa:", len(h.clients))
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
				SenderID: uint(senderID),
			}
			_, err = h.messageService.SaveMessage(dbMessage)
			if err != nil {
				log.Printf("Gagal menyimpan pesan ke database: %v", err)
				// Kita tetap melanjutkan broadcast meskipun gagal menyimpan
			}

			// Marshal pesan ke JSON untuk logging
			if msgJson, err := json.Marshal(chatMessage); err == nil {
				log.Printf("Broadcasting pesan: %s", msgJson)
			}

			for client := range h.clients {
				select {
				case client.Send <- chatMessage:
				default:
					// Jika channel send penuh, client dianggap lambat.
					// Kita tutup koneksinya dan hapus dari hub.
					close(client.Send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// sendChatHistory mengambil riwayat chat dari database dan mengirimkannya ke satu client.
func (h *Hub) sendChatHistory(client *Client) {
	history, err := h.messageService.GetMessage()
	if err != nil {
		log.Printf("Gagal mengambil riwayat chat: %v", err)
		return
	}

	log.Printf("Mengirim %d pesan riwayat ke client %s", len(history), client.UserID)

	for _, msg := range history {
		chatMessage := ChatMessage{
			Type:      "history_message", // Tipe pesan khusus untuk riwayat
			Content:   msg.Content,
			SenderID:  strconv.FormatUint(uint64(msg.SenderID), 10),
			Timestamp: msg.CreatedAt,
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
