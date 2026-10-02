package routes

import (
	"notification-service/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	notifications := router.Group("/notifications")
	{
		notifications.POST("", handlers.CreateNotification)
		notifications.GET("/:email", handlers.GetNotificationsByEmail)
		notifications.PATCH("/:id/read", handlers.MarkAsRead)
		notifications.GET("/unread-count/:email", handlers.GetUnreadCount)
	}

	certificates := router.Group("/certificates")
	{
		certificates.POST("", handlers.GenerateCertificate)
		certificates.GET("/:id/download", handlers.DownloadCertificate)
		certificates.GET("/student/:email", handlers.GetCertificatesByEmail)
	}
	analytics := router.Group("/analytics")
	{
		analytics.GET("/platform", handlers.GetPlatformStats)
		analytics.GET("/certificates-by-course", handlers.GetCertificatesByCourse)
	}
}