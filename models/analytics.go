package models

type PlatformStats struct {
	TotalNotifications  int `json:"totalNotifications"`
	UnreadNotifications int `json:"unreadNotifications"`
	TotalCertificates   int `json:"totalCertificates"`
}

type CertificateStats struct {
	CourseTitle      string `json:"courseTitle"`
	CertificateCount int    `json:"certificateCount"`
}