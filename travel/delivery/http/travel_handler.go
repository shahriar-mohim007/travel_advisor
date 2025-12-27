package http

import (
	"encoding/json"
	"net/http"
	"travel_advisor/domain"
	"travel_advisor/helpers"
	"travel_advisor/travel/transformer"

	"github.com/go-chi/chi/v5"
)

type TravelHandler struct {
	TravelUsecase domain.TravelUsecase
}

func NewTravelHandler(r *chi.Mux, t domain.TravelUsecase) {
	handler := &TravelHandler{
		TravelUsecase: t,
	}
	r.Route("/v1/travel", func(r chi.Router) {
		r.Use(helpers.JWTAuthMiddleware)
		r.Get("/coolest/districts", handler.List)
		r.Post("/recommend", handler.Recommend)
	})
}

func (h *TravelHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	districts, err := h.TravelUsecase.CoolestDistricts(ctx)
	if err != nil {
		resp := &helpers.Response{
			Status:  http.StatusInternalServerError,
			Message: "coolest districts fetch failed",
			Error:   err.Error(),
		}
		resp.Render(w)
		return
	}

	resp := &helpers.Response{
		Status: http.StatusOK,
		Data:   transformer.TransformCoolestDistrictResponse(districts),
	}
	resp.Render(w)
}

func (h *TravelHandler) Recommend(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req domain.TravelRecommendationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp := &helpers.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		}
		resp.Render(w)
		return
	}

	result, err := h.TravelUsecase.RecommendTravel(ctx, req)
	if err != nil {
		resp := &helpers.Response{
			Status:  http.StatusInternalServerError,
			Message: "Travel recommendation failed",
			Error:   err.Error(),
		}
		resp.Render(w)
		return
	}

	resp := &helpers.Response{
		Status: http.StatusOK,
		Data:   result,
	}
	resp.Render(w)
}
