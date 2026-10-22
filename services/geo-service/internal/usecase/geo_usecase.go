package usecase

import (
	"context"
	"errors"

	"sakid/geo-service/internal/domain"
)

type geoUseCase struct {
	repo domain.GeoRepository
}

func NewGeoUseCase(repo domain.GeoRepository) domain.GeoUseCase {
	return &geoUseCase{repo: repo}
}

func (u *geoUseCase) ResolveLocation(ctx context.Context, input domain.ReverseGeocodeInput) (*domain.ReverseGeocodeResponse, error) {
	if input.Latitude == 0 && input.Longitude == 0 {
		return nil, errors.New("พิกัด GPS ไม่ถูกต้อง")
	}
	return u.repo.FindAuthorityByLocation(ctx, input.Latitude, input.Longitude)
}