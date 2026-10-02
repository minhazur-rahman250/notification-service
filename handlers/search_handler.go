package handlers

import (
	"net/http"
	"notification-service/config"
	"notification-service/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// IndexCourse - POST /search/index
// NestJS থেকে course তৈরি/publish/update হলে এটা কল হবে, search index আপডেট রাখার জন্য
func IndexCourse(c *gin.Context) {
	var req models.IndexCourseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ON CONFLICT - একই course_id আগে থেকে থাকলে নতুন করে insert না করে update করে দেয়
	query := `
		INSERT INTO course_search_index (course_id, title, description, price, teacher_name)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (course_id)
		DO UPDATE SET title = $2, description = $3, price = $4, teacher_name = $5, indexed_at = CURRENT_TIMESTAMP
		RETURNING id, course_id, title, description, price, teacher_name, indexed_at`

	var course models.CourseIndex
	err := config.DB.QueryRow(
		query, req.CourseID, req.Title, req.Description, req.Price, req.TeacherName,
	).Scan(
		&course.ID, &course.CourseID, &course.Title, &course.Description,
		&course.Price, &course.TeacherName, &course.IndexedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "course index করতে ব্যর্থ"})
		return
	}

	c.JSON(http.StatusOK, course)
}

// SearchCourses - GET /search/courses?q=nestjs&minPrice=0&maxPrice=100&sort=price_asc
func SearchCourses(c *gin.Context) {
	searchQuery := c.Query("q")
	minPriceStr := c.DefaultQuery("minPrice", "0")
	maxPriceStr := c.DefaultQuery("maxPrice", "999999")
	sort := c.DefaultQuery("sort", "newest")

	minPrice, err := strconv.ParseFloat(minPriceStr, 64)
	if err != nil {
		minPrice = 0
	}
	maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
	if err != nil {
		maxPrice = 999999
	}

	orderBy := "indexed_at DESC"
	switch sort {
	case "price_asc":
		orderBy = "price ASC"
	case "price_desc":
		orderBy = "price DESC"
	case "title_asc":
		orderBy = "title ASC"
	}

	// ILIKE - PostgreSQL এ case-insensitive pattern matching (বড়/ছোট হাতের অক্ষর নিয়ে ভাবতে হয় না)
	query := `
		SELECT id, course_id, title, description, price, teacher_name, indexed_at
		FROM course_search_index
		WHERE (title ILIKE $1 OR description ILIKE $1)
		AND price BETWEEN $2 AND $3
		ORDER BY ` + orderBy

	rows, err := config.DB.Query(query, "%"+searchQuery+"%", minPrice, maxPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "সার্চ করতে ব্যর্থ"})
		return
	}
	defer rows.Close()

	results := []models.CourseIndex{}
	for rows.Next() {
		var course models.CourseIndex
		if err := rows.Scan(
			&course.ID, &course.CourseID, &course.Title, &course.Description,
			&course.Price, &course.TeacherName, &course.IndexedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ডেটা পড়তে ব্যর্থ"})
			return
		}
		results = append(results, course)
	}

	c.JSON(http.StatusOK, gin.H{
		"results": results,
		"count":   len(results),
	})
}

// RemoveCourseFromIndex - DELETE /search/index/:courseId
// NestJS এ course ডিলিট হলে index থেকেও সরাতে হবে
func RemoveCourseFromIndex(c *gin.Context) {
	courseId := c.Param("courseId")

	_, err := config.DB.Exec(`DELETE FROM course_search_index WHERE course_id = $1`, courseId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "index থেকে সরাতে ব্যর্থ"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "index থেকে সরানো হয়েছে"})
}