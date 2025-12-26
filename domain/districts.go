package domain

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
