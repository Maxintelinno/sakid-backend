package domain

import (
	"context"
	"time"
)

type ReportStatus string

const (
	StatusPending    ReportStatus = "pending"     // รอดำเนินการ (สีส้ม)
	StatusInProgress ReportStatus = "in_progress" // กำลังแก้ไข (สีน้ำเงิน)
	StatusResolved   ReportStatus = "resolved"    // แก้ไขเสร็จสิ้น (สีเขียว)
	StatusRejected   ReportStatus = "rejected"    // ไม่อยู่ในอำนาจ/ปฏิเสธ
)

type Report struct {
	ID             string       `json:"id"`
	CaseCode       string       `json:"case_code"` // รหัสเคส เช่น CV-2569-012
	UserID         string       `json:"user_id"`
	Title          string       `json:"title"`
	Description    string       `json:"description"`
	CategoryID     string       `json:"category_id"`
	CategoryName   string       `json:"category_name"`
	Latitude       float64      `json:"latitude"`
	Longitude      float64      `json:"longitude"`
	AddressHint    string       `json:"address_hint"`
	MediaURLs      []string     `json:"media_urls"`
	IsAnonymous    bool         `json:"is_anonymous"`
	AutoBlurMedia  bool         `json:"auto_blur_media"`
	Status         ReportStatus `json:"status"`
	UpvoteCount    int          `json:"upvote_count"`
	TargetAgencyID string       `json:"target_agency_id,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

type CreateReportInput struct {
	UserID        string   `json:"-"`
	Title         string   `json:"title" validate:"required"`
	Description   string   `json:"description" validate:"required"`
	CategoryID    string   `json:"category_id" validate:"required"`
	Latitude      float64  `json:"latitude" validate:"required"`
	Longitude     float64  `json:"longitude" validate:"required"`
	AddressHint   string   `json:"address_hint"`
	MediaURLs     []string `json:"media_urls"`
	IsAnonymous   bool     `json:"is_anonymous"`
	AutoBlurMedia bool     `json:"auto_blur_media"`
}

// โครงสร้างสำหรับระบบแจ้งเตือนปัญหาซ้ำ (Duplicate Alert)
type CheckDuplicateInput struct {
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	CategoryID string  `json:"category_id"`
	RadiusM    float64 `json:"radius_meters"` // รัศมีค้นหา เช่น 50 เมตร
}

type DuplicateReportSummary struct {
	ID          string  `json:"id"`
	CaseCode    string  `json:"case_code"`
	Title       string  `json:"title"`
	DistanceM   float64 `json:"distance_meters"`
	Status      string  `json:"status"`
	UpvoteCount int     `json:"upvote_count"`
	Thumbnail   string  `json:"thumbnail,omitempty"`
}

// Interfaces
type ReportRepository interface {
	Create(ctx context.Context, report *Report) error
	GetByID(ctx context.Context, id string) (*Report, error)
	List(ctx context.Context, limit, offset int) ([]Report, error)
	FindDuplicates(ctx context.Context, input CheckDuplicateInput) ([]DuplicateReportSummary, error)
	Upvote(ctx context.Context, reportID string, userID string) (int, error)
}

type ReportUseCase interface {
	CreateReport(ctx context.Context, input CreateReportInput) (*Report, error)
	GetReport(ctx context.Context, id string) (*Report, error)
	ListReports(ctx context.Context, limit, offset int) ([]Report, error)
	CheckDuplicates(ctx context.Context, input CheckDuplicateInput) ([]DuplicateReportSummary, error)
	UpvoteReport(ctx context.Context, reportID string, userID string) (int, error)
}