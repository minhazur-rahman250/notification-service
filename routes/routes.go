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
		
	}
	certificates := router.Group("/certificates")
	{
		certificates.POST("", handlers.GenerateCertificate)
		certificates.GET("/:id/download", handlers.DownloadCertificate)
		certificates.GET("/student/:email", handlers.GetCertificatesByEmail)
	}
}