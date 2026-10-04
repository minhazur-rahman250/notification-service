package handlers

import (
	"net/http"
	"notification-service/config"
	"notification-service/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// IndexCourse - POST /search/index
// NestJS থেকে course create/update হলে এটা কল হয়, Course+User+Lesson+Enrollment
// এর সারসংক্ষেপ এখানে সংরক্ষিত থাকে যাতে সার্চ দ্রুত হয়, মূল টেবিলে বারবার জয়েন লাগে না
func IndexCourse(c *gin.Context) {
	var req models.IndexCourseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
		INSERT INTO course_search_index
			(course_id, title, description, price, teacher_name, lesson_count, enrollment_count, is_published)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (course_id)
		DO UPDATE SET
			title = $2, description = $3, price = $4, teacher_name = $5,
			lesson_count = $6, enrollment_count = $7, is_published = $8,
			indexed_at = CURRENT_TIMESTAMP
		RETURNING id, course_id, title, description, price, teacher_name, lesson_count, enrollment_count, is_published, indexed_at`

	var course models.CourseIndex
	err := config.DB.QueryRow(
		query, req.CourseID, req.Title, req.Description, req.Price, req.TeacherName,
		req.LessonCount, req.EnrollmentCount, req.IsPublished,
	).Scan(
		&course.ID, &course.CourseID, &course.Title, &course.Description, &course.Price,
		&course.TeacherName, &course.LessonCount, &course.EnrollmentCount,
		&course.IsPublished, &course.IndexedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "course index করতে ব্যর্থ"})
		return
	}

	c.JSON(http.StatusOK, course)
}

// SearchCourses - GET /search/courses?q=&minPrice=&maxPrice=&sort=
// sort: newest (default) | price_asc | price_desc | title_asc | popular
// শুধু is_published = true কোর্স পাবলিক সার্চে দেখায়
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
	case "popular":
		orderBy = "enrollment_count DESC"
	}

	query := `
		SELECT id, course_id, title, description, price, teacher_name, lesson_count, enrollment_count, is_published, indexed_at
		FROM course_search_index
		WHERE (title ILIKE $1 OR description ILIKE $1)
		AND price BETWEEN $2 AND $3
		AND is_published = TRUE
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
			&course.ID, &course.CourseID, &course.Title, &course.Description, &course.Price,
			&course.TeacherName, &course.LessonCount, &course.EnrollmentCount,
			&course.IsPublished, &course.IndexedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ডেটা পড়তে ব্যর্থ"})
			return
		}
		results = append(results, course)
	}

	c.JSON(http.StatusOK, gin.H{"results": results, "count": len(results)})
}

// RemoveCourseFromIndex - DELETE /search/index/:courseId
func RemoveCourseFromIndex(c *gin.Context) {
	courseId := c.Param("courseId")

	_, err := config.DB.Exec(`DELETE FROM course_search_index WHERE course_id = $1`, courseId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "index থেকে সরাতে ব্যর্থ"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "index থেকে সরানো হয়েছে"})
}