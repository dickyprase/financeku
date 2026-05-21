package handler

import (
	"encoding/json"
	"net/http"

	"github.com/financeku/backend/internal/service"
	"github.com/financeku/backend/pkg/response"
	"github.com/financeku/backend/pkg/validator"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input service.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate
	errors := validator.ValidationErrors{}
	validator.Required(errors, "name", input.Name)
	validator.Required(errors, "email", input.Email)
	validator.Email(errors, "email", input.Email)
	validator.Required(errors, "password", input.Password)
	validator.MinLength(errors, "password", input.Password, 6)

	if errors.HasErrors() {
		response.ValidationError(w, errors)
		return
	}

	user, err := h.authService.Register(input)
	if err != nil {
		response.Error(w, http.StatusConflict, err.Error())
		return
	}

	response.Created(w, "Registration successful", user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input service.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate
	errors := validator.ValidationErrors{}
	validator.Required(errors, "email", input.Email)
	validator.Required(errors, "password", input.Password)

	if errors.HasErrors() {
		response.ValidationError(w, errors)
		return
	}

	user, tokens, err := h.authService.Login(input)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.OK(w, "Login successful", map[string]interface{}{
		"user":   user,
		"tokens": tokens,
	})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if body.RefreshToken == "" {
		response.Error(w, http.StatusBadRequest, "Refresh token is required")
		return
	}

	tokens, err := h.authService.RefreshToken(body.RefreshToken)
	if err != nil {
		response.Unauthorized(w, err.Error())
		return
	}

	response.OK(w, "Token refreshed", tokens)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// In a stateless JWT setup, logout is handled client-side
	// Optionally implement token blacklisting here
	response.OK(w, "Logged out successfully", nil)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	user, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		response.NotFound(w, "User not found")
		return
	}
	response.OK(w, "User profile", user)
}
