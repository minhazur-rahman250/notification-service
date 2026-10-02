package handlers

import (
	"fmt"
	"net/http"
	"notification-service/config"
	"notification-service/models"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
)

// GenerateCertificate - POST /certificates
func GenerateCertificate(c *gin.Context) {
	var req models.GenerateCertificateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// প্রতিটা certificate এর জন্য একটা unique নম্বর (recruiter/employer verify করতে পারবে)
	certNumber := fmt.Sprintf("CERT-%s", uuid.New().String()[:8])
	fileName := fmt.Sprintf("%s.pdf", certNumber)
	filePath := fmt.Sprintf("./certificates/%s", fileName)

	// certificates ফোল্ডার না থাকলে বানিয়ে নেয়
	if err := os.MkdirAll("./certificates", os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ফোল্ডার তৈরি করতে ব্যর্থ"})
		return
	}

	if err := createPDF(filePath, req.StudentName, req.CourseTitle, certNumber); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "PDF তৈরি করতে ব্যর্থ"})
		return
	}

	var certificate models.Certificate
	query := `
		INSERT INTO certificates (certificate_number, student_email, student_name, course_title, file_path)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, certificate_number, student_email, student_name, course_title, file_path, issued_at`

	err := config.DB.QueryRow(
		query, certNumber, req.StudentEmail, req.StudentName, req.CourseTitle, filePath,
	).Scan(
		&certificate.ID, &certificate.CertificateNumber, &certificate.StudentEmail,
		&certificate.StudentName, &certificate.CourseTitle, &certificate.FilePath, &certificate.IssuedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "certificate রেকর্ড সেভ করতে ব্যর্থ"})
		return
	}

	c.JSON(http.StatusCreated, certificate)
}

// createPDF - gofpdf দিয়ে আসল certificate ডিজাইন করা হচ্ছে
func createPDF(filePath, studentName, courseTitle, certNumber string) error {
	pdf := gofpdf.New("L", "mm", "A4", "") // L = Landscape, A4 সাইজ
	pdf.AddPage()

	// বর্ডার আঁকা
	pdf.SetLineWidth(1.5)
	pdf.Rect(10, 10, 277, 190, "D")

	pdf.SetFont("Arial", "B", 28)
	pdf.SetY(40)
	pdf.CellFormat(297, 15, "Certificate of Completion", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 14)
	pdf.SetY(70)
	pdf.CellFormat(297, 10, "This is to certify that", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "B", 22)
	pdf.SetY(85)
	pdf.CellFormat(297, 12, studentName, "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 14)
	pdf.SetY(105)
	pdf.CellFormat(297, 10, "has successfully completed the course", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "B", 18)
	pdf.SetY(118)
	pdf.CellFormat(297, 10, courseTitle, "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.SetY(160)
	pdf.CellFormat(297, 8, fmt.Sprintf("Certificate No: %s", certNumber), "", 1, "C", false, 0, "")
	pdf.CellFormat(297, 8, fmt.Sprintf("Issued on: %s", time.Now().Format("January 2, 2006")), "", 1, "C", false, 0, "")

	return pdf.OutputFileAndClose(filePath)
}

// DownloadCertificate - GET /certificates/:id/download
func DownloadCertificate(c *gin.Context) {
	id := c.Param("id")

	var filePath, certNumber string
	query := `SELECT file_path, certificate_number FROM certificates WHERE id = $1`

	err := config.DB.QueryRow(query, id).Scan(&filePath, &certNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "certificate পাওয়া যায়নি"})
		return
	}

	// ব্রাউজারে সরাসরি PDF download শুরু করিয়ে দেয়
	c.FileAttachment(filePath, certNumber+".pdf")
}

// GetCertificatesByEmail - GET /certificates/student/:email
func GetCertificatesByEmail(c *gin.Context) {
	email := c.Param("email")

	query := `
		SELECT id, certificate_number, student_email, student_name, course_title, file_path, issued_at
		FROM certificates WHERE student_email = $1 ORDER BY issued_at DESC`

	rows, err := config.DB.Query(query, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "certificate খুঁজতে ব্যর্থ"})
		return
	}
	defer rows.Close()

	certificates := []models.Certificate{}
	for rows.Next() {
		var cert models.Certificate
		if err := rows.Scan(
			&cert.ID, &cert.CertificateNumber, &cert.StudentEmail,
			&cert.StudentName, &cert.CourseTitle, &cert.FilePath, &cert.IssuedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ডেটা পড়তে ব্যর্থ"})
			return
		}
		certificates = append(certificates, cert)
	}

	c.JSON(http.StatusOK, certificates)
}