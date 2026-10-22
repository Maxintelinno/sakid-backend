package repository

import (
	"context"
	"fmt"
	"math"

	"sakid/geo-service/internal/domain"
)

type geoRepository struct {
	authorities []domain.Authority
}

func NewGeoRepository() domain.GeoRepository {
	// ตัวอย่างข้อมูล Master Data ของ อปท. (ในระบบจริงดึงจากตาราง PostGIS)
	data := []domain.Authority{
		{
			ID:            "bma-chatuchak",
			Name:          "ฝ่ายโยธา สำนักงานเขตจตุจักร",
			Type:          domain.TypeBMA,
			Province:      "กรุงเทพมหานคร",
			District:      "เขตจตุจักร",
			Subdistrict:   "แขวงจันทรเกษม",
			PostalCode:    "10900",
			OfficialEmail: "saraban.chatuchak@bangkok.go.th",
			BackupEmail:   "damrongdhama.bkk@moi.go.th",
			PhoneNumber:   "02-513-3444",
		},
		{
			ID:            "sao-kudhae",
			Name:          "องค์การบริหารส่วนตำบลกุดแห่",
			Type:          domain.TypeSAO,
			Province:      "สกลนคร",
			District:      "อำเภอเลิงนกทา",
			Subdistrict:   "ตำบลกุดแห่",
			PostalCode:    "47120",
			OfficialEmail: "saraban@kudhae.go.th",
			BackupEmail:   "damrongdhama.skn@moi.go.th",
			PhoneNumber:   "042-789-111",
		},
	}
	return &geoRepository{authorities: data}
}

func (r *geoRepository) FindAuthorityByLocation(ctx context.Context, lat, lon float64) (*domain.ReverseGeocodeResponse, error) {
	// ในระบบจริงใช้คำสั่ง SQL:
	// SELECT * FROM authorities WHERE ST_Contains(geom, ST_SetSRID(ST_Point(lon, lat), 4326)) LIMIT 1;
	
	// Mock Fallback ตรวจจับพิกัดกว้างๆ ใน กทม.
	if math.Abs(lat-13.8) < 0.5 && math.Abs(lon-100.5) < 0.5 {
		auth := r.authorities[0]
		return &domain.ReverseGeocodeResponse{
			Subdistrict: auth.Subdistrict,
			District:    auth.District,
			Province:    auth.Province,
			PostalCode:  auth.PostalCode,
			Formatted:   fmt.Sprintf("%s %s %s %s", auth.Subdistrict, auth.District, auth.Province, auth.PostalCode),
			Authority:   auth,
		}, nil
	}

	// Default กรณีอยู่นอกพื้นที่จำลอง
	auth := r.authorities[1]
	return &domain.ReverseGeocodeResponse{
		Subdistrict: auth.Subdistrict,
		District:    auth.District,
		Province:    auth.Province,
		PostalCode:  auth.PostalCode,
		Formatted:   fmt.Sprintf("%s %s %s %s", auth.Subdistrict, auth.District, auth.Province, auth.PostalCode),
		Authority:   auth,
	}, nil
}

func (r *geoRepository) ListAuthorities(ctx context.Context, province string) ([]domain.Authority, error) {
	return r.authorities, nil
}