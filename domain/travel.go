package domain

import "context"

type TravelRecommendationRequest struct {
	CurrentLat          float64 `json:"current_lat"`
	CurrentLong         float64 `json:"current_long"`
	DestinationDistrict string  `json:"destination_district"`
	TravelDate          string  `json:"travel_date"`
}

type TravelRecommendationResponse struct {
	Destination    string  `json:"destination"`
	Recommendation string  `json:"recommendation"`
	Reason         string  `json:"reason"`
	TempDiff       float64 `json:"temp_diff"`
	PM25Diff       float64 `json:"pm25_diff"`
}

type TravelUsecase interface {
	CoolestDistricts(ctx context.Context) ([]DistrictCache, error)
	RecommendTravel(ctx context.Context, req TravelRecommendationRequest) (*TravelRecommendationResponse, error)
}
