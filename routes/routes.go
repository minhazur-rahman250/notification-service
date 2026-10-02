package routes

import (
	"notification-service/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/notifications")
	{
		api.POST("", handlers.CreateNotification)
		api.GET("/:email", handlers.GetNotificationsByEmail)
		api.PATCH("/:id/read", handlers.MarkAsRead)             // নতুন
		api.GET("/unread-count/:email", handlers.GetUnreadCount) // নতুন
	}
	

	
}