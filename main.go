package main

import (
	"example/hello/book"
	"example/hello/realtime"
	"example/hello/route"
	"example/hello/short"
	"example/hello/user"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Muat variabel dari file .env
	err := godotenv.Load()
	if err != nil {
		// Jangan hentikan aplikasi jika .env tidak ada, karena bisa di-set di environment produksi
		log.Println("Peringatan: Gagal memuat file .env")
	}

	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connected to the database successfully")

	db.AutoMigrate(&book.Book{})
	db.AutoMigrate(&user.User{})
	db.AutoMigrate(&short.Short{})
	db.AutoMigrate(&realtime.Message{})

	// Inisialisasi komponen untuk Message
	messageRepository := realtime.NewRepository(db)
	messageService := realtime.NewService(messageRepository)

	// Buat dan jalankan Hub real-time dalam goroutine terpisah
	hub := realtime.NewHub(messageService)
	go hub.Run()

	// Create a new Gin router
	r := gin.Default()

	route.BookRoutes(r, db)
	route.UserRoutes(r, db)
	route.ShortRoutes(r, db)
	route.WebSocketRoutes(r, hub)

	// Start the server on port 8080
	r.Run(":8080")
}

// main
// handler
// service
// repository
// request
// db
// mysql

// Bagaimana cara membuat "room" atau "channel" agar pesan tidak di-broadcast ke semua orang?
