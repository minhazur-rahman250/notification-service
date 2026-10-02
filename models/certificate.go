package models

import "time"

type Certificate struct {
	ID                int       `json:"id"`
	CertificateNumber string    `json:"certificateNumber"`
	StudentEmail      string    `json:"studentEmail"`
	StudentName       string    `json:"studentName"`
	CourseTitle       string    `json:"courseTitle"`
	FilePath          string    `json:"filePath"`
	IssuedAt          time.Time `json:"issuedAt"`
}

// NestJS থেকে এই shape-এ ডেটা পাঠাতে হবে
type GenerateCertificateRequest struct {
	StudentEmail string `json:"studentEmail" binding:"required,email"`
	StudentName  string `json:"studentName" binding:"required"`
	CourseTitle  string `json:"courseTitle" binding:"required"`
}