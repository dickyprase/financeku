package handler

import (
	"encoding/json"
	"net/http"

	"github.com/financeku/backend/internal/service"
	"github.com/financeku/backend/pkg/response"
	"github.com/financeku/backend/pkg/validator"
)

type GoalHandler struct {
	goalService *service.GoalService
}

func NewGoalHandler(goalService *service.GoalService) *GoalHandler {
	return &GoalHandler{goalService: goalService}
}

func (h *GoalHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)

	var input service.GoalInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	errors := validator.ValidationErrors{}
	validator.Required(errors, "name", input.Name)
	validator.MinValue(errors, "target_amount", input.TargetAmount, 1)
	validator.InList(errors, "tracking_mode", input.TrackingMode, []string{"manual", "single_wallet", "multiple_wallet", "all_wallet"})

	if errors.HasErrors() {
		response.ValidationError(w, errors)
		return
	}

	goal, err := h.goalService.Create(userID, input)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.Created(w, "Goal created", goal)
}

func (h *GoalHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)

	goals, err := h.goalService.List(userID)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.OK(w, "Goals", goals)
}

func (h *GoalHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	id := getPathParam(r, "id")

	goal, err := h.goalService.GetByID(userID, id)
	if err != nil {
		response.NotFound(w, "Goal not found")
		return
	}

	response.OK(w, "Goal", goal)
}

func (h *GoalHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	id := getPathParam(r, "id")

	var input service.GoalInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	errors := validator.ValidationErrors{}
	validator.Required(errors, "name", input.Name)
	validator.MinValue(errors, "target_amount", input.TargetAmount, 1)

	if errors.HasErrors() {
		response.ValidationError(w, errors)
		return
	}

	goal, err := h.goalService.Update(userID, id, input)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, "Goal updated", goal)
}

func (h *GoalHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	id := getPathParam(r, "id")

	if err := h.goalService.Delete(userID, id); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, "Goal deleted", nil)
}

func (h *GoalHandler) GetProgress(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	id := getPathParam(r, "id")

	progress, err := h.goalService.GetProgress(userID, id)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, "Goal progress", progress)
}
