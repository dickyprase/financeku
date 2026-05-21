package handler

import (
	"encoding/json"
	"net/http"

	"github.com/financeku/backend/internal/service"
	"github.com/financeku/backend/pkg/response"
	"github.com/financeku/backend/pkg/validator"
)

type ProfileHandler struct {
	profileService *service.ProfileService
}

func NewProfileHandler(profileService *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

func (h *ProfileHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)

	var input service.ProfileUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	errors := validator.ValidationErrors{}
	validator.Required(errors, "name", input.Name)

	if errors.HasErrors() {
		response.ValidationError(w, errors)
		return
	}

	user, err := h.profileService.Update(userID, input)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, "Profile updated", user)
}

func (h *ProfileHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)

	var input service.ChangePasswordInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	errors := validator.ValidationErrors{}
	validator.Required(errors, "old_password", input.OldPassword)
	validator.Required(errors, "new_password", input.NewPassword)
	validator.MinLength(errors, "new_password", input.NewPassword, 6)

	if errors.HasErrors() {
		response.ValidationError(w, errors)
		return
	}

	if err := h.profileService.ChangePassword(userID, input); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, "Password changed successfully", nil)
}

// Admin Handler

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page := getQueryInt(r, "page", 1)
	perPage := getQueryInt(r, "per_page", 20)

	users, total, err := h.adminService.ListUsers(page, perPage)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	response.SuccessWithMeta(w, http.StatusOK, "Users", users, &response.Meta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

func (h *AdminHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input service.AdminCreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

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

	user, err := h.adminService.CreateUser(input)
	if err != nil {
		response.Error(w, http.StatusConflict, err.Error())
		return
	}

	response.Created(w, "User created", user)
}

func (h *AdminHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := getPathParam(r, "id")

	var input service.AdminCreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.adminService.UpdateUser(id, input)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, "User updated", user)
}

func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := getPathParam(r, "id")

	if err := h.adminService.DeleteUser(id); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, "User deleted", nil)
}

func (h *AdminHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	id := getPathParam(r, "id")

	var body struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	errors := validator.ValidationErrors{}
	validator.Required(errors, "password", body.Password)
	validator.MinLength(errors, "password", body.Password, 6)

	if errors.HasErrors() {
		response.ValidationError(w, errors)
		return
	}

	if err := h.adminService.ResetPassword(id, body.Password); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, "Password reset successful", nil)
}

func (h *AdminHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.adminService.GetSettings()
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.OK(w, "Site settings", settings)
}

func (h *AdminHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	var body map[string]string
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	for key, value := range body {
		if err := h.adminService.UpdateSetting(key, value); err != nil {
			response.Error(w, http.StatusBadRequest, "Failed to update setting: "+key)
			return
		}
	}

	response.OK(w, "Settings updated", nil)
}

func (h *AdminHandler) GetActivityLogs(w http.ResponseWriter, r *http.Request) {
	page := getQueryInt(r, "page", 1)
	perPage := getQueryInt(r, "per_page", 50)

	logs, total, err := h.adminService.GetActivityLogs(page, perPage)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	response.SuccessWithMeta(w, http.StatusOK, "Activity logs", logs, &response.Meta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}
