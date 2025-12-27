package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"
	"travel_advisor/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations
type MockCache struct {
	mock.Mock
}

func (m *MockCache) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockCache) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockCache) Keys(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

type MockDistrictRepository struct {
	mock.Mock
}

func (m *MockDistrictRepository) List(ctx context.Context, ctr *domain.DistrictCriteria) ([]*domain.District, error) {
	args := m.Called(ctx, ctr)
	return args.Get(0).([]*domain.District), args.Error(1)
}

func TestTravelUsecase_CoolestDistricts(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*MockCache, *MockDistrictRepository)
		expectedResult []domain.DistrictCache
		expectedError  error
	}{
		{
			name: "Success - Returns sorted districts",
			setupMocks: func(mockCache *MockCache, mockDistrictRepo *MockDistrictRepository) {

				mockCache.On("Keys", mock.Anything).Return([]string{"Dhaka", "Chittagong", "Sylhet"}, nil)

				dhakaData := domain.DistrictCache{Name: "Dhaka", AvgTemp2PM: 30.5, AvgPM25: 45.2}
				chittagongData := domain.DistrictCache{Name: "Chittagong", AvgTemp2PM: 28.3, AvgPM25: 35.1}
				sylhetData := domain.DistrictCache{Name: "Sylhet", AvgTemp2PM: 26.8, AvgPM25: 25.5}

				dhakaJSON, _ := json.Marshal(dhakaData)
				chittagongJSON, _ := json.Marshal(chittagongData)
				sylhetJSON, _ := json.Marshal(sylhetData)

				mockCache.On("Get", mock.Anything, "Dhaka").Return(string(dhakaJSON), nil)
				mockCache.On("Get", mock.Anything, "Chittagong").Return(string(chittagongJSON), nil)
				mockCache.On("Get", mock.Anything, "Sylhet").Return(string(sylhetJSON), nil)
			},
			expectedResult: []domain.DistrictCache{
				{Name: "Sylhet", AvgTemp2PM: 26.8, AvgPM25: 25.5},
				{Name: "Chittagong", AvgTemp2PM: 28.3, AvgPM25: 35.1},
				{Name: "Dhaka", AvgTemp2PM: 30.5, AvgPM25: 45.2},
			},
			expectedError: nil,
		},
		{
			name: "Error - Cache keys failure",
			setupMocks: func(mockCache *MockCache, mockDistrictRepo *MockDistrictRepository) {
				mockCache.On("Keys", mock.Anything).Return([]string{}, errors.New("cache error"))
			},
			expectedResult: nil,
			expectedError:  errors.New("cache error"),
		},
		{
			name: "Success - Empty cache",
			setupMocks: func(mockCache *MockCache, mockDistrictRepo *MockDistrictRepository) {
				mockCache.On("Keys", mock.Anything).Return([]string{}, nil)
			},
			expectedResult: nil,
			expectedError:  nil,
		},
		{
			name: "Success - More than 10 districts returns top 10",
			setupMocks: func(mockCache *MockCache, mockDistrictRepo *MockDistrictRepository) {

				districts := make([]string, 12)
				for i := 0; i < 12; i++ {
					districts[i] = fmt.Sprintf("District%d", i)
				}
				mockCache.On("Keys", mock.Anything).Return(districts, nil)

				for i := 0; i < 12; i++ {
					districtData := domain.DistrictCache{
						Name:       fmt.Sprintf("District%d", i),
						AvgTemp2PM: float64(20 + i),
						AvgPM25:    float64(10 + i),
					}
					districtJSON, _ := json.Marshal(districtData)
					mockCache.On("Get", mock.Anything, fmt.Sprintf("District%d", i)).Return(string(districtJSON), nil)
				}
			},
			expectedResult: func() []domain.DistrictCache {
				result := make([]domain.DistrictCache, 10)
				for i := 0; i < 10; i++ {
					result[i] = domain.DistrictCache{
						Name:       fmt.Sprintf("District%d", i),
						AvgTemp2PM: float64(20 + i),
						AvgPM25:    float64(10 + i),
					}
				}
				return result
			}(),
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := new(MockCache)
			mockDistrictRepo := new(MockDistrictRepository)

			tt.setupMocks(mockCache, mockDistrictRepo)

			usecase := NewTravelUsecase(mockCache, mockDistrictRepo)

			result, err := usecase.CoolestDistricts(context.Background())

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			mockCache.AssertExpectations(t)
			mockDistrictRepo.AssertExpectations(t)
		})
	}
}

func TestTravelUsecase_RecommendTravel(t *testing.T) {
	tests := []struct {
		name           string
		request        domain.TravelRecommendationRequest
		setupMocks     func(*MockCache, *MockDistrictRepository)
		expectedResult *domain.TravelRecommendationResponse
		expectedError  error
	}{
		{
			name: "Success - Recommended travel",
			request: domain.TravelRecommendationRequest{
				CurrentLat:          23.7104,
				CurrentLong:         90.3944,
				DestinationDistrict: "Sylhet",
				TravelDate:          "2024-01-15",
			},
			setupMocks: func(mockCache *MockCache, mockDistrictRepo *MockDistrictRepository) {
				district := &domain.District{
					ID:   1,
					Name: "Sylhet",
					Lat:  24.8949,
					Long: 91.8687,
				}
				mockDistrictRepo.On("List", mock.Anything, &domain.DistrictCriteria{
					DistrictName: stringPtr("Sylhet"),
				}).Return([]*domain.District{district}, nil)
			},
			expectedResult: &domain.TravelRecommendationResponse{
				Destination:    "Sylhet",
				Recommendation: "Recommended",
				Reason:         "Your destination is 2.0°C cooler and has significantly better air quality. Enjoy your trip!",
				TempDiff:       -2.0,
				PM25Diff:       -5.0,
			},
			expectedError: nil,
		},
		{
			name: "Error - District not found",
			request: domain.TravelRecommendationRequest{
				CurrentLat:          23.7104,
				CurrentLong:         90.3944,
				DestinationDistrict: "NonExistent",
				TravelDate:          "2024-01-15",
			},
			setupMocks: func(mockCache *MockCache, mockDistrictRepo *MockDistrictRepository) {
				mockDistrictRepo.On("List", mock.Anything, &domain.DistrictCriteria{
					DistrictName: stringPtr("NonExistent"),
				}).Return([]*domain.District{}, errors.New("district not found"))
			},
			expectedResult: nil,
			expectedError:  errors.New("destination district not found"),
		},
		{
			name: "Success - Not recommended travel",
			request: domain.TravelRecommendationRequest{
				CurrentLat:          23.7104,
				CurrentLong:         90.3944,
				DestinationDistrict: "Dhaka",
				TravelDate:          "2024-01-15",
			},
			setupMocks: func(mockCache *MockCache, mockDistrictRepo *MockDistrictRepository) {
				district := &domain.District{
					ID:   2,
					Name: "Dhaka",
					Lat:  23.8103,
					Long: 90.4125,
				}
				mockDistrictRepo.On("List", mock.Anything, &domain.DistrictCriteria{
					DistrictName: stringPtr("Dhaka"),
				}).Return([]*domain.District{district}, nil)
			},
			expectedResult: &domain.TravelRecommendationResponse{
				Destination:    "Dhaka",
				Recommendation: "Not Recommended",
				Reason:         "Your destination is hotter and has worse air quality than your current location. It's better to stay where you are.",
				TempDiff:       3.5,
				PM25Diff:       10.2,
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := new(MockCache)
			mockDistrictRepo := new(MockDistrictRepository)

			tt.setupMocks(mockCache, mockDistrictRepo)

			usecase := NewTravelUsecase(mockCache, mockDistrictRepo)

			result, err := usecase.RecommendTravel(context.Background(), tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				if result != nil {
					assert.Equal(t, tt.expectedResult.Destination, result.Destination)
				}
			}

			mockCache.AssertExpectations(t)
			mockDistrictRepo.AssertExpectations(t)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
