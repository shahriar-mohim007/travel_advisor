package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"travel_advisor/domain"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock TravelUsecase
type MockTravelUsecase struct {
	mock.Mock
}

func (m *MockTravelUsecase) CoolestDistricts(ctx context.Context) ([]domain.DistrictCache, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.DistrictCache), args.Error(1)
}

func (m *MockTravelUsecase) RecommendTravel(ctx context.Context, req domain.TravelRecommendationRequest) (*domain.TravelRecommendationResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TravelRecommendationResponse), args.Error(1)
}

func TestTravelHandler_List(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*MockTravelUsecase)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success - Returns coolest districts",
			setupMocks: func(mockUsecase *MockTravelUsecase) {
				districts := []domain.DistrictCache{
					{Name: "Sylhet", AvgTemp2PM: 26.8, AvgPM25: 25.5},
					{Name: "Chittagong", AvgTemp2PM: 28.3, AvgPM25: 35.1},
				}
				mockUsecase.On("CoolestDistricts", mock.Anything).Return(districts, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"data":`,
		},
		{
			name: "Error - Usecase returns error",
			setupMocks: func(mockUsecase *MockTravelUsecase) {
				mockUsecase.On("CoolestDistricts", mock.Anything).Return([]domain.DistrictCache{}, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"message":"coolest districts fetch failed"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockTravelUsecase)
			tt.setupMocks(mockUsecase)

			handler := &TravelHandler{
				TravelUsecase: mockUsecase,
			}

			req := httptest.NewRequest("GET", "/v1/travel/coolest/districts", nil)
			rr := httptest.NewRecorder()

			handler.List(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestTravelHandler_Recommend(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMocks     func(*MockTravelUsecase)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success - Returns travel recommendation",
			requestBody: domain.TravelRecommendationRequest{
				CurrentLat:          23.7104,
				CurrentLong:         90.3944,
				DestinationDistrict: "Sylhet",
				TravelDate:          "2024-01-15",
			},
			setupMocks: func(mockUsecase *MockTravelUsecase) {
				req := domain.TravelRecommendationRequest{
					CurrentLat:          23.7104,
					CurrentLong:         90.3944,
					DestinationDistrict: "Sylhet",
					TravelDate:          "2024-01-15",
				}
				resp := &domain.TravelRecommendationResponse{
					Destination:    "Sylhet",
					Recommendation: "Recommended",
					Reason:         "Great weather!",
					TempDiff:       -2.5,
					PM25Diff:       -5.0,
				}
				mockUsecase.On("RecommendTravel", mock.Anything, req).Return(resp, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"data":`,
		},
		{
			name:        "Error - Invalid JSON body",
			requestBody: "invalid json",
			setupMocks: func(mockUsecase *MockTravelUsecase) {

			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"message":"Invalid request body"`,
		},
		{
			name: "Error - Usecase returns error",
			requestBody: domain.TravelRecommendationRequest{
				CurrentLat:          23.7104,
				CurrentLong:         90.3944,
				DestinationDistrict: "NonExistent",
				TravelDate:          "2024-01-15",
			},
			setupMocks: func(mockUsecase *MockTravelUsecase) {
				req := domain.TravelRecommendationRequest{
					CurrentLat:          23.7104,
					CurrentLong:         90.3944,
					DestinationDistrict: "NonExistent",
					TravelDate:          "2024-01-15",
				}
				mockUsecase.On("RecommendTravel", mock.Anything, req).Return(nil, errors.New("district not found"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"message":"Travel recommendation failed"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(MockTravelUsecase)
			tt.setupMocks(mockUsecase)

			handler := &TravelHandler{
				TravelUsecase: mockUsecase,
			}

			var body bytes.Buffer
			if tt.requestBody != nil {
				if str, ok := tt.requestBody.(string); ok {
					body.WriteString(str)
				} else {
					json.NewEncoder(&body).Encode(tt.requestBody)
				}
			}

			req := httptest.NewRequest("POST", "/v1/travel/recommend", &body)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handler.Recommend(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestNewTravelHandler(t *testing.T) {
	r := chi.NewRouter()
	mockUsecase := new(MockTravelUsecase)

	assert.NotPanics(t, func() {
		NewTravelHandler(r, mockUsecase)
	})
}
