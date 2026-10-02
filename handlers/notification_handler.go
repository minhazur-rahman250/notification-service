package handlers

import (
	"net/http"
	"notification-service/config"
	"notification-service/models"

	"github.com/gin-gonic/gin"
)

// CreateNotification - POST /notifications
// NestJS থেকে কল হবে যখন কোনো event ঘটবে (enrollment, course publish ইত্যাদি)
func CreateNotification(c *gin.Context) {
	var req models.CreateNotificationRequest

	// c.ShouldBindJSON - request body পড়ে req struct এ বসায়, সাথে validation করে
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var notification models.Notification
	query := `
		INSERT INTO notifications (recipient_email, type, message)
		VALUES ($1, $2, $3)
		RETURNING id, recipient_email, type, message, is_read, created_at`

	err := config.DB.QueryRow(
		query, req.RecipientEmail, req.Type, req.Message,
	).Scan(
		&notification.ID, &notification.RecipientEmail, &notification.Type,
		&notification.Message, &notification.IsRead, &notification.CreatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "notification সেভ করতে ব্যর্থ"})
		return
	}

	c.JSON(http.StatusCreated, notification)
}

// GetNotificationsByEmail - GET /notifications/:email
func GetNotificationsByEmail(c *gin.Context) {
	email := c.Param("email")

	query := `
		SELECT id, recipient_email, type, message, is_read, created_at
		FROM notifications
		WHERE recipient_email = $1
		ORDER BY created_at DESC`

	rows, err := config.DB.Query(query, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "notification খুঁজতে ব্যর্থ"})
		return
	}
	defer rows.Close() // ফাংশন শেষ হওয়ার সাথে সাথে rows বন্ধ করে দেয় - memory leak আটকায়

	notifications := []models.Notification{}
	for rows.Next() {
		var n models.Notification
		if err := rows.Scan(&n.ID, &n.RecipientEmail, &n.Type, &n.Message, &n.IsRead, &n.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ডেটা পড়তে ব্যর্থ"})
			return
		}
		notifications = append(notifications, n)
	}

	c.JSON(http.StatusOK, notifications)
}
// MarkAsRead - PATCH /notifications/:id/read
func MarkAsRead(c *gin.Context) {
	id := c.Param("id")

	query := `UPDATE notifications SET is_read = TRUE WHERE id = $1 RETURNING id, recipient_email, type, message, is_read, created_at`

	var notification models.Notification
	err := config.DB.QueryRow(query, id).Scan(
		&notification.ID, &notification.RecipientEmail, &notification.Type,
		&notification.Message, &notification.IsRead, &notification.CreatedAt,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification পাওয়া যায়নি"})
		return
	}

	c.JSON(http.StatusOK, notification)
}

// GetUnreadCount - GET /notifications/unread-count/:email
func GetUnreadCount(c *gin.Context) {
	email := c.Param("email")

	var count int
	query := `SELECT COUNT(*) FROM notifications WHERE recipient_email = $1 AND is_read = FALSE`

	err := config.DB.QueryRow(query, email).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "count বের করতে ব্যর্থ"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unreadCount": count})
}