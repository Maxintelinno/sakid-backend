package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	appHTTP "sakid/geo-service/internal/delivery/http"
	"sakid/geo-service/internal/repository"
	"sakid/geo-service/internal/usecase"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083" // พอร์ต 8083 สำหรับ geo-service
	}

	repo := repository.NewGeoRepository()
	uc := usecase.NewGeoUseCase(repo)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, echo.Map{"status": "ok", "service": "geo-service"})
	})

	appHTTP.RegisterGeoRoutes(e, uc)

	log.Printf("🚀 Geo Service กำลังทำงานบนพอร์ต :%s\n", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Server ล้มเหลว: %v", err)
	}
}