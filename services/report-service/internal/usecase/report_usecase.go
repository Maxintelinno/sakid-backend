package usecase

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"sakid/report-service/internal/domain"

	"github.com/google/uuid"
)

type reportUseCase struct {
	repo domain.ReportRepository
}

func NewReportUseCase(repo domain.ReportRepository) domain.ReportUseCase {
	return &reportUseCase{repo: repo}
}

func (u *reportUseCase) CreateReport(ctx context.Context, input domain.CreateReportInput) (*domain.Report, error) {
	yearBE := time.Now().Year() + 543 // แปลง พ.ศ.
	caseCode := fmt.Sprintf("CV-%d-%03d", yearBE, rand.Intn(900)+100)

	report := &domain.Report{
		ID:            uuid.New().String(),
		CaseCode:      caseCode,
		UserID:        input.UserID,
		Title:         input.Title,
		Description:   input.Description,
		CategoryID:    input.CategoryID,
		Latitude:      input.Latitude,
		Longitude:     input.Longitude,
		AddressHint:   input.AddressHint,
		MediaURLs:     input.MediaURLs,
		IsAnonymous:   input.IsAnonymous,
		AutoBlurMedia: input.AutoBlurMedia,
		Status:        domain.StatusPending,
		UpvoteCount:   1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err := u.repo.Create(ctx, report)
	if err != nil {
		return nil, err
	}

	// TODO: ส่ง Event ไปยัง Message Queue (Redis / Kafka) เพื่อให้ dispatcher-worker นำไปยิง Email หาท้องถิ่น

	return report, nil
}

func (u *reportUseCase) GetReport(ctx context.Context, id string) (*domain.Report, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *reportUseCase) ListReports(ctx context.Context, limit, offset int) ([]domain.Report, error) {
	return u.repo.List(ctx, limit, offset)
}

func (u *reportUseCase) CheckDuplicates(ctx context.Context, input domain.CheckDuplicateInput) ([]domain.DuplicateReportSummary, error) {
	if input.RadiusM == 0 {
		input.RadiusM = 50 // ค่าเริ่มต้นตรวจจับในรัศมี 50 เมตร
	}
	return u.repo.FindDuplicates(ctx, input)
}

func (u *reportUseCase) UpvoteReport(ctx context.Context, reportID string, userID string) (int, error) {
	return u.repo.Upvote(ctx, reportID, userID)
}