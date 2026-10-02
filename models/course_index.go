package models

import "time"

type CourseIndex struct {
	ID          int       `json:"id"`
	CourseID    int       `json:"courseId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	TeacherName string    `json:"teacherName"`
	IndexedAt   time.Time `json:"indexedAt"`
}

// NestJS থেকে course তৈরি/আপডেট হলে এই shape-এ ডেটা পাঠাবে
type IndexCourseRequest struct {
	CourseID    int     `json:"courseId" binding:"required"`
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	TeacherName string  `json:"teacherName"`
}