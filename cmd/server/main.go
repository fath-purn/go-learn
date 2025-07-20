package main

import (
	"example/hello/internal/book"
	"example/hello/internal/handler"
	"example/hello/internal/realtime"
	"example/hello/internal/route"
	"example/hello/internal/short"
	"example/hello/internal/user"
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

	// === Dependency Injection Setup ===
	// Inisialisasi semua dependency di satu tempat (Composition Root)

	// User Dependencies
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	// Book Dependencies
	bookRepository := book.NewRepository(db)
	bookService := book.NewService(bookRepository)
	bookHandler := handler.NewBookHandler(bookService)

	// Short URL Dependencies
	shortRepository := short.NewRepository(db)
	shortService := short.NewService(shortRepository)
	shortHandler := handler.NewShortUrlHandler(shortService)

	// Buat dan jalankan Hub real-time dalam goroutine terpisah
	messageRepository := realtime.NewRepository(db)
	messageService := realtime.NewService(messageRepository)
	hub := realtime.NewHub(messageService, userService)
	go hub.Run()
	webSocketHandler := handler.NewWebSocketHandler(hub)

	// Create a new Gin router
	r := gin.Default()

	// Setup routes dengan menyuntikkan handler yang sudah dibuat
	route.SetupRoutes(r, userHandler, bookHandler, shortHandler, webSocketHandler)

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
