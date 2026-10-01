package models

import "time"

// এই struct টা ঠিক NestJS এর Entity-র মতো কাজ করে - ডেটার shape define করে
type Notification struct {
	ID             int       `json:"id"`
	RecipientEmail string    `json:"recipientEmail"`
	Type           string    `json:"type"`
	Message        string    `json:"message"`
	IsRead         bool      `json:"isRead"`
	CreatedAt      time.Time `json:"createdAt"`
}

// এটা incoming request এর body validate/bind করার জন্য (NestJS এর DTO এর সমতুল্য)
type CreateNotificationRequest struct {
	RecipientEmail string `json:"recipientEmail" binding:"required,email"`
	Type           string `json:"type" binding:"required"`
	Message        string `json:"message" binding:"required"`
}