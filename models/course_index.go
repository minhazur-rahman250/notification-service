package models

import "time"

type CourseIndex struct {
	ID              int       `json:"id"`
	CourseID        int       `json:"courseId"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Price           float64   `json:"price"`
	TeacherName     string    `json:"teacherName"`
	LessonCount     int       `json:"lessonCount"`
	EnrollmentCount int       `json:"enrollmentCount"`
	IsPublished     bool      `json:"isPublished"`
	IndexedAt       time.Time `json:"indexedAt"`
}

type IndexCourseRequest struct {
	CourseID        int     `json:"courseId" binding:"required"`
	Title           string  `json:"title" binding:"required"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	TeacherName     string  `json:"teacherName"`
	LessonCount     int     `json:"lessonCount"`
	EnrollmentCount int     `json:"enrollmentCount"`
	IsPublished     bool    `json:"isPublished"`
}