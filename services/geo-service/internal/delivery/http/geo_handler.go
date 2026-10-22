package http

import (
	"net/http"

	"sakid/geo-service/internal/domain"

	"github.com/labstack/echo/v4"
)

type GeoHandler struct {
	usecase domain.GeoUseCase
}

func RegisterGeoRoutes(e *echo.Echo, uc domain.GeoUseCase) {
	h := &GeoHandler{usecase: uc}

	v1 := e.Group("/api/v1/geo")
	v1.POST("/reverse-geocode", h.ReverseGeocode)
}

func (h *GeoHandler) ReverseGeocode(c echo.Context) error {
	var input domain.ReverseGeocodeInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "ข้อมูลพิกัดไม่ถูกต้อง"})
	}

	result, err := h.usecase.ResolveLocation(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}