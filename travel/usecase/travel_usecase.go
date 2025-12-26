package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"travel_advisor/domain"
	"travel_advisor/helpers"
	"travel_advisor/pkg/cache"
	"travel_advisor/pkg/conn"
	"travel_advisor/pkg/log"
)

type TravelUsecase struct {
	CacheRepository     cache.Cache
	DistrictsRepository domain.DistrictRepository
}

func NewTravelUsecase(c cache.Cache, d domain.DistrictRepository) domain.TravelUsecase {
	return &TravelUsecase{
		CacheRepository:     c,
		DistrictsRepository: d,
	}
}

func (t *TravelUsecase) CoolestDistricts(ctx context.Context) ([]domain.DistrictCache, error) {

	districtNames, err := t.CacheRepository.Keys(ctx)
	log.Println(districtNames)
	log.Println(len(districtNames))

	if err != nil {
		return nil, err
	}
	var districts []domain.DistrictCache

	for _, name := range districtNames {
		dataStr, err := t.CacheRepository.Get(ctx, name)
		if err != nil {
			continue
		}

		var d domain.DistrictCache
		if err := json.Unmarshal([]byte(dataStr), &d); err != nil {
			continue
		}

		districts = append(districts, d)
	}

	sort.Slice(districts, func(i, j int) bool {
		if districts[i].AvgTemp2PM == districts[j].AvgTemp2PM {
			return districts[i].AvgPM25 < districts[j].AvgPM25
		}
		return districts[i].AvgTemp2PM < districts[j].AvgTemp2PM
	})

	if len(districts) > 10 {
		districts = districts[:10]
	}

	return districts, nil
}

func (t *TravelUsecase) RecommendTravel(
	ctx context.Context,
	req domain.TravelRecommendationRequest,
) (*domain.TravelRecommendationResponse, error) {

	districts, err := t.DistrictsRepository.List(ctx, &domain.DistrictCriteria{
		DistrictName: &req.DestinationDistrict,
	})
	if err != nil || len(districts) == 0 {
		return nil, fmt.Errorf("destination district not found")
	}
	destDistrict := districts[0]

	conn.InitClient()
	client := conn.GetHTTClient()

	date := req.TravelDate

	var (
		destTemp    float64
		destPM25    float64
		currentTemp float64
		currentPM25 float64

		errDestTemp error
		errDestPM25 error
		errCurTemp  error
		errCurPM25  error
	)

	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		defer wg.Done()
		destTemp, errDestTemp = helpers.FetchAvgTempAt2PM(
			ctx,
			client,
			destDistrict.Lat,
			destDistrict.Long,
			&date,
		)
	}()

	go func() {
		defer wg.Done()
		destPM25, errDestPM25 = helpers.FetchAvgPM25(
			ctx,
			client,
			destDistrict.Lat,
			destDistrict.Long,
			&date,
		)
	}()

	go func() {
		defer wg.Done()
		currentTemp, errCurTemp = helpers.FetchAvgTempAt2PM(
			ctx,
			client,
			req.CurrentLat,
			req.CurrentLong,
			&date,
		)
	}()

	go func() {
		defer wg.Done()
		currentPM25, errCurPM25 = helpers.FetchAvgPM25(
			ctx,
			client,
			req.CurrentLat,
			req.CurrentLong,
			&date,
		)
	}()

	wg.Wait()

	if errDestTemp != nil {
		return nil, errDestTemp
	}
	if errDestPM25 != nil {
		return nil, errDestPM25
	}
	if errCurTemp != nil {
		return nil, errCurTemp
	}
	if errCurPM25 != nil {
		return nil, errCurPM25
	}

	tempDiff := destTemp - currentTemp
	pm25Diff := destPM25 - currentPM25

	resp := &domain.TravelRecommendationResponse{
		Destination:    destDistrict.Name,
		TempDiff:       tempDiff,
		PM25Diff:       pm25Diff,
		Recommendation: "Not Recommended",
	}

	if tempDiff < 0 && pm25Diff < 0 {
		resp.Recommendation = "Recommended"
		resp.Reason = fmt.Sprintf(
			"Your destination is %.1f°C cooler and has significantly better air quality. Enjoy your trip!",
			-tempDiff,
		)
	} else {
		resp.Reason = "Your destination is hotter and has worse air quality than your current location." +
			" It’s better to stay where you are."
	}

	return resp, nil
}
