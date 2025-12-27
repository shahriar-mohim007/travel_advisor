package http

import (
	"encoding/json"
	"net/http"
	"travel_advisor/domain"
	"travel_advisor/helpers"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	UserUsecase domain.UserUsecase
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func NewUserHandler(r *chi.Mux, u domain.UserUsecase) {
	handler := &UserHandler{
		UserUsecase: u,
	}

	r.Route("/v1/auth", func(r chi.Router) {
		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)
	})
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp := &helpers.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		}
		resp.Render(w)
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		resp := &helpers.Response{
			Status:  http.StatusBadRequest,
			Message: "Name, email, and password are required",
		}
		resp.Render(w)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		resp := &helpers.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to hash password",
			Error:   err.Error(),
		}
		resp.Render(w)
		return
	}

	user := &domain.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = h.UserUsecase.Create(ctx, user)
	if err != nil {
		resp := &helpers.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create user",
			Error:   err.Error(),
		}
		resp.Render(w)
		return
	}

	resp := &helpers.Response{
		Status:  http.StatusCreated,
		Message: "User registered successfully",
	}
	resp.Render(w)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp := &helpers.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Error:   err.Error(),
		}
		resp.Render(w)
		return
	}

	if req.Email == "" || req.Password == "" {
		resp := &helpers.Response{
			Status:  http.StatusBadRequest,
			Message: "Email and password are required",
		}
		resp.Render(w)
		return
	}

	user, err := h.UserUsecase.Get(ctx, &domain.UserCriteria{
		Email: &req.Email,
	})
	if err != nil {
		resp := &helpers.Response{
			Status:  http.StatusUnauthorized,
			Message: "Invalid credentials",
		}
		resp.Render(w)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		resp := &helpers.Response{
			Status:  http.StatusUnauthorized,
			Message: "Invalid credentials",
		}
		resp.Render(w)
		return
	}

	token, err := helpers.GenerateJWTToken(user.ID)
	if err != nil {
		resp := &helpers.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to generate token",
			Error:   err.Error(),
		}
		resp.Render(w)
		return
	}

	authResponse := &AuthResponse{
		Token: token,
	}

	resp := &helpers.Response{
		Status:  http.StatusOK,
		Message: "Login successful",
		Data:    authResponse,
	}
	resp.Render(w)
}
