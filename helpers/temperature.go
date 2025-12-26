package helpers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"travel_advisor/pkg/log"
)

func FetchAvgTempAt2PM(
	ctx context.Context,
	client *http.Client,
	lat, long float64, date *string,
) (float64, error) {

	params := url.Values{}
	params.Set("latitude", fmt.Sprintf("%f", lat))
	params.Set("longitude", fmt.Sprintf("%f", long))
	params.Set("hourly", "temperature_2m")
	params.Set("timezone", "auto")
	if date != nil {
		params.Set("start_date", *date)
		params.Set("end_date", *date)
	} else {
		params.Set("forecast_days", "7")
	}

	weatherURL := "https://api.open-meteo.com/v1/forecast?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, weatherURL, nil)
	if err != nil {
		return 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data struct {
		Hourly struct {
			Temperature []float64 `json:"temperature_2m"`
		} `json:"hourly"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	temps := data.Hourly.Temperature
	log.Println(temps)

	var sum float64
	var count int
	for i := 14; i < len(temps); i += 24 {
		sum += temps[i]
		count++
	}

	if count == 0 {
		return 0, errors.New("no temperature data")
	}

	return sum / float64(count), nil
}

func FetchAvgPM25(
	ctx context.Context,
	client *http.Client,
	lat, long float64, date *string,
) (float64, error) {

	params := url.Values{}
	params.Set("latitude", fmt.Sprintf("%f", lat))
	params.Set("longitude", fmt.Sprintf("%f", long))
	params.Set("hourly", "pm2_5")
	params.Set("timezone", "auto")

	if date != nil {
		params.Set("start_date", *date)
		params.Set("end_date", *date)
	} else {
		params.Set("forecast_days", "7")
	}

	airURL := "https://air-quality-api.open-meteo.com/v1/air-quality?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, airURL, nil)
	if err != nil {
		return 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data struct {
		Hourly struct {
			PM25 []float64 `json:"pm2_5"`
		} `json:"hourly"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	log.Println(data.Hourly.PM25)
	var sum float64
	for _, v := range data.Hourly.PM25 {
		sum += v
	}

	if len(data.Hourly.PM25) == 0 {
		return 0, errors.New("no PM2.5 data")
	}

	return sum / float64(len(data.Hourly.PM25)), nil
}
