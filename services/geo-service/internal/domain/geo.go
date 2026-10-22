package domain

import "context"

type AuthorityType string

const (
	TypeBMA        AuthorityType = "bma_district" // สำนักงานเขต กทม.
	TypeCityMunicipality  AuthorityType = "city_municipality" // เทศบาลนคร
	TypeTownMunicipality  AuthorityType = "town_municipality" // เทศบาลเมือง
	TypeSubdistrictMunicipality AuthorityType = "subdistrict_municipality" // เทศบาลตำบล
	TypeSAO        AuthorityType = "sao" // องค์การบริหารส่วนตำบล (อบต.)
)

type Authority struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`           // เช่น สำนักงานเขตจตุจักร, เทศบาลตำบลกุดแห่
	Type          AuthorityType `json:"type"`
	Province      string        `json:"province"`
	District      string        `json:"district"`
	Subdistrict   string        `json:"subdistrict"`
	PostalCode    string        `json:"postal_code"`
	OfficialEmail string        `json:"official_email"` // อีเมลสารบรรณกลาง
	BackupEmail   string        `json:"backup_email"`   // อีเมลศูนย์ดำรงธรรม
	PhoneNumber   string        `json:"phone_number"`
}

type ReverseGeocodeInput struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

type ReverseGeocodeResponse struct {
	Subdistrict string    `json:"subdistrict"`
	District    string    `json:"district"`
	Province    string    `json:"province"`
	PostalCode  string    `json:"postal_code"`
	Formatted   string    `json:"formatted_address"`
	Authority   Authority `json:"responsible_authority"`
}

type GeoRepository interface {
	FindAuthorityByLocation(ctx context.Context, lat, lon float64) (*ReverseGeocodeResponse, error)
	ListAuthorities(ctx context.Context, province string) ([]Authority, error)
}

type GeoUseCase interface {
	ResolveLocation(ctx context.Context, input ReverseGeocodeInput) (*ReverseGeocodeResponse, error)
}