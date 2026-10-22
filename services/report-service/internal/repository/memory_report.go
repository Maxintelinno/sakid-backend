package repository

import (
	"context"
	"errors"
	"math"
	"sync"

	"sakid/report-service/internal/domain"
)

type memoryReportRepo struct {
	mu      sync.RWMutex
	reports map[string]domain.Report
}

func NewMemoryReportRepository() domain.ReportRepository {
	return &memoryReportRepo{
		reports: make(map[string]domain.Report),
	}
}

func (r *memoryReportRepo) Create(ctx context.Context, report *domain.Report) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reports[report.ID] = *report
	return nil
}

func (r *memoryReportRepo) GetByID(ctx context.Context, id string) (*domain.Report, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	report, exists := r.reports[id]
	if !exists {
		return nil, errors.New("ไม่พบข้อมูลเรื่องร้องเรียน")
	}
	return &report, nil
}

func (r *memoryReportRepo) List(ctx context.Context, limit, offset int) ([]domain.Report, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]domain.Report, 0, len(r.reports))
	for _, item := range r.reports {
		list = append(list, item)
	}
	return list, nil
}

// ตรวจจับปัญหาซ้ำด้วยสูตรคำนวณระยะทาง Haversine Formula (เมตร)
func (r *memoryReportRepo) FindDuplicates(ctx context.Context, input domain.CheckDuplicateInput) ([]domain.DuplicateReportSummary, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var duplicates []domain.DuplicateReportSummary
	for _, report := range r.reports {
		if report.Status == domain.StatusResolved {
			continue // ข้ามเรื่องที่ซ่อมเสร็จแล้ว
		}

		dist := calculateDistanceMeters(input.Latitude, input.Longitude, report.Latitude, report.Longitude)
		if dist <= input.RadiusM {
			thumb := ""
			if len(report.MediaURLs) > 0 {
				thumb = report.MediaURLs[0]
			}
			duplicates = append(duplicates, domain.DuplicateReportSummary{
				ID:          report.ID,
				CaseCode:    report.CaseCode,
				Title:       report.Title,
				DistanceM:   math.Round(dist),
				Status:      string(report.Status),
				UpvoteCount: report.UpvoteCount,
				Thumbnail:   thumb,
			})
		}
	}
	return duplicates, nil
}

func (r *memoryReportRepo) Upvote(ctx context.Context, reportID string, userID string) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	report, exists := r.reports[reportID]
	if !exists {
		return 0, errors.New("ไม่พบเรื่องร้องเรียน")
	}
	report.UpvoteCount++
	r.reports[reportID] = report
	return report.UpvoteCount, nil
}

func calculateDistanceMeters(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000 // เมตร
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*(math.Pi/180.0))*math.Cos(lat2*(math.Pi/180.0))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}