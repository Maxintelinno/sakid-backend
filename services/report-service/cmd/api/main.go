package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	appHTTP "sakid/report-service/internal/delivery/http"
	"sakid/report-service/internal/repository"
	"sakid/report-service/internal/usecase"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082" // พอร์ต 8082 แยกจาก auth-service (8081)
	}

	repo := repository.NewMemoryReportRepository()
	uc := usecase.NewReportUseCase(repo)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, echo.Map{"status": "ok", "service": "report-service"})
	})

	appHTTP.RegisterReportRoutes(e, uc)

	log.Printf("🚀 Report Service กำลังทำงานบนพอร์ต :%s\n", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Server ล้มเหลว: %v", err)
	}
}