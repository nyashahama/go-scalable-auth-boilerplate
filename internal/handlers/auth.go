package handlers

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"time"

	"user-auth-app/internal/services"
	"user-auth-app/internal/validator"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type AuthHandler struct {
	service *services.AuthService
	logger  *zerolog.Logger
	timeout time.Duration
}

func NewAuthHandler(service *services.AuthService, logger *zerolog.Logger, timeout time.Duration) *AuthHandler {
	return &AuthHandler{service: service, logger: logger, timeout: timeout}
}

type ErrorResponse struct {
	Error  string                      `json:"error"`
	Fields []validator.ValidationError `json:"fields,omitempty"`
}

func (h *AuthHandler) writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func (h *AuthHandler) writeValidationError(w http.ResponseWriter, errors []validator.ValidationError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:  "validation failed",
		Fields: errors,
	})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	v := validator.New()
	v.ValidateUsername("username", req.Username)
	v.ValidateEmail("email", req.Email)
	v.ValidatePassword("password", req.Password)
	v.ValidateRole("role", req.Role)

	if !v.Valid() {
		h.writeValidationError(w, v.Errors())
		return
	}

	// Default role
	if req.Role == "" {
		req.Role = "user"
	}

	user, err := h.service.Register(ctx, req.Username, req.Email, req.Password, req.Role)
	if err != nil {
		h.logger.Error().Err(err).Msg("registration failed")
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	v := validator.New()
	v.ValidateRequired("email", req.Email)
	v.ValidateRequired("password", req.Password)

	if !v.Valid() {
		h.writeValidationError(w, v.Errors())
		return
	}

	token, err := h.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		h.writeError(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	userIDStr := chi.URLParam(r, "id")
	userID64, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		h.writeError(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	if userID64 < math.MinInt32 || userID64 > math.MaxInt32 {
		h.writeError(w, "user ID out of range", http.StatusBadRequest)
		return
	}

	userID := int32(userID64)

	user, err := h.service.GetProfile(ctx, userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get profile")
		h.writeError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if user.ID == 0 {
		h.writeError(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
