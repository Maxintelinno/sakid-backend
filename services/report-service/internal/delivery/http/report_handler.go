package http

import (
	"net/http"
	"strconv"

	"sakid/report-service/internal/domain"

	"github.com/labstack/echo/v4"
)

type ReportHandler struct {
	usecase domain.ReportUseCase
}

func RegisterReportRoutes(e *echo.Echo, uc domain.ReportUseCase) {
	h := &ReportHandler{usecase: uc}

	v1 := e.Group("/api/v1/reports")
	v1.GET("", h.ListReports)
	v1.POST("", h.CreateReport)
	v1.GET("/:id", h.GetReport)
	v1.POST("/check-duplicates", h.CheckDuplicates)
	v1.POST("/:id/upvote", h.UpvoteReport)
}

func (h *ReportHandler) CreateReport(c echo.Context) error {
	var input domain.CreateReportInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "ข้อมูลนำเข้าไม่ถูกต้อง"})
	}

	// ในระบบจริง User ID จะถูกดึงมาจาก JWT Token Middleware
	input.UserID = "usr_demo_123"

	report, err := h.usecase.CreateReport(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, report)
}

func (h *ReportHandler) CheckDuplicates(c echo.Context) error {
	var input domain.CheckDuplicateInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "ข้อมูลพิกัดไม่ถูกต้อง"})
	}

	duplicates, err := h.usecase.CheckDuplicates(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"has_duplicates": len(duplicates) > 0,
		"items":          duplicates,
	})
}

func (h *ReportHandler) ListReports(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if limit <= 0 {
		limit = 20
	}

	reports, err := h.usecase.ListReports(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"data": reports})
}

func (h *ReportHandler) GetReport(c echo.Context) error {
	id := c.Param("id")
	report, err := h.usecase.GetReport(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, report)
}

func (h *ReportHandler) UpvoteReport(c echo.Context) error {
	id := c.Param("id")
	newCount, err := h.usecase.UpvoteReport(c.Request().Context(), id, "usr_demo_123")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"upvote_count": newCount})
}