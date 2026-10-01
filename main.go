package main

import (
	"log"
	"os"

	"notification-service/config"
	"notification-service/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// .env ফাইল লোড করা - না পেলেও crash করাচ্ছি না (production এ env vars সরাসরি set থাকতে পারে)
	if err := godotenv.Load(); err != nil {
		log.Println(".env ফাইল পাওয়া যায়নি, সিস্টেম env variable ব্যবহার হচ্ছে")
	}

	config.ConnectDatabase()

	router := gin.Default() // Default() built-in logger আর crash-recovery middleware সহ আসে
	routes.SetupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Notification service চালু হচ্ছে পোর্ট %s এ\n", port)
	router.Run(":" + port)
}