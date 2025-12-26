package domain

import (
	"context"
	"errors"
	"time"
)

type DistrictResponse struct {
	Districts []DistrictJSON `json:"districts"`
}

type DistrictJSON struct {
	Id         string `json:"id"`
	DivisionID string `json:"division_id"`
	Name       string `json:"name"`
	BnName     string `json:"bn_name"`
	Lat        string `json:"lat"`
	Long       string `json:"long"`
}

type DistrictCache struct {
	Name       string
	AvgTemp2PM float64
	AvgPM25    float64
}
type DistrictCriteria struct {
	DistrictName *string
}
type District struct {
	ID         int64     `json:"id"`
	DivisionID int       `json:"division_id"`
	Name       string    `json:"name"`
	BnName     string    `json:"bn_name"`
	Lat        float64   `json:"lat"`
	Long       float64   `json:"long"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type DistrictRepository interface {
	List(ctx context.Context, ctr *DistrictCriteria) ([]*District, error)
}

var (
	ErrDistrictNotFound = errors.New("district not found")
)
