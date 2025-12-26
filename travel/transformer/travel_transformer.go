package transformer

import "travel_advisor/domain"

type DistrictResponse struct {
	Name       string  `json:"district_name"`
	AvgTemp2PM float64 `json:"avg_temp_2_pm"`
	AvgPM25    float64 `json:"avg_pm_25"`
}

func TransformCoolestDistrictResponse(gt []domain.DistrictCache) []DistrictResponse {
	resp := make([]DistrictResponse, 0)
	for _, g := range gt {
		resp = append(resp, DistrictResponse{
			Name:       g.Name,
			AvgTemp2PM: g.AvgTemp2PM,
			AvgPM25:    g.AvgPM25,
		})
	}
	return resp
}
