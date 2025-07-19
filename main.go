package main

import (
	"example/hello/book"
	"example/hello/route"
	"example/hello/short"
	"example/hello/user"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/djawa?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connected to the database successfully")

	db.AutoMigrate(&book.Book{})
	db.AutoMigrate(&user.User{})
	db.AutoMigrate(&short.Short{})

	// Create a new Gin router
	r := gin.Default()

	route.BookRoutes(r, db)
	route.UserRoutes(r, db)
	route.ShortRoutes(r, db)

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
