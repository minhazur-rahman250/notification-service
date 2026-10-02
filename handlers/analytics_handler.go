package handlers

import (
	"net/http"
	"notification-service/config"
	"notification-service/models"
	"sync"

	"github.com/gin-gonic/gin"
)

// GetPlatformStats - GET /analytics/platform
// তিনটা আলাদা COUNT query একসাথে (goroutine দিয়ে) parallel-এ চালানো হচ্ছে,
// একটার পর একটা sequential চালালে যতটা সময় লাগত তার চেয়ে দ্রুত শেষ হয়
func GetPlatformStats(c *gin.Context) {
	var stats models.PlatformStats
	var wg sync.WaitGroup
	var mu sync.Mutex // একসাথে stats struct-এ লেখার সময় data race আটকাতে
	errChan := make(chan error, 3)

	wg.Add(3)

	go func() {
		defer wg.Done()
		var count int
		err := config.DB.QueryRow(`SELECT COUNT(*) FROM notifications`).Scan(&count)
		if err != nil {
			errChan <- err
			return
		}
		mu.Lock()
		stats.TotalNotifications = count
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		var count int
		err := config.DB.QueryRow(`SELECT COUNT(*) FROM notifications WHERE is_read = FALSE`).Scan(&count)
		if err != nil {
			errChan <- err
			return
		}
		mu.Lock()
		stats.UnreadNotifications = count
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		var count int
		err := config.DB.QueryRow(`SELECT COUNT(*) FROM certificates`).Scan(&count)
		if err != nil {
			errChan <- err
			return
		}
		mu.Lock()
		stats.TotalCertificates = count
		mu.Unlock()
	}()

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "পরিসংখ্যান বের করতে ব্যর্থ"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetCertificatesByCourse - GET /analytics/certificates-by-course
// কোন course-এ সবচেয়ে বেশি student সার্টিফিকেট পেয়েছে, সেই র‍্যাংকিং
func GetCertificatesByCourse(c *gin.Context) {
	query := `
		SELECT course_title, COUNT(*) as certificate_count
		FROM certificates
		GROUP BY course_title
		ORDER BY certificate_count DESC`

	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "পরিসংখ্যান বের করতে ব্যর্থ"})
		return
	}
	defer rows.Close()

	results := []models.CertificateStats{}
	for rows.Next() {
		var stat models.CertificateStats
		if err := rows.Scan(&stat.CourseTitle, &stat.CertificateCount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ডেটা পড়তে ব্যর্থ"})
			return
		}
		results = append(results, stat)
	}

	c.JSON(http.StatusOK, results)
}