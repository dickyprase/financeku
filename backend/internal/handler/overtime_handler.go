package handler

import (
	"encoding/json"
	"net/http"

	"github.com/financeku/backend/internal/service"
	"github.com/financeku/backend/pkg/response"
	"github.com/financeku/backend/pkg/validator"
)

type OvertimeHandler struct {
	overtimeService *service.OvertimeService
}

func NewOvertimeHandler(overtimeService *service.OvertimeService) *OvertimeHandler {
	return &OvertimeHandler{overtimeService: overtimeService}
}

func (h *OvertimeHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)

	var input service.OvertimeInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	errors := validator.ValidationErrors{}
	validator.Required(errors, "date", input.Date)
	validator.MinValue(errors, "hours", input.Hours, 0.5)
	validator.MaxValue(errors, "hours", input.Hours, 12)

	if errors.HasErrors() {
		response.ValidationError(w, errors)
		return
	}

	record, err := h.overtimeService.Create(userID, input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Created(w, "Overtime record created", record)
}

func (h *OvertimeHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	month := getQueryString(r, "month", "")
	page := getQueryInt(r, "page", 1)
	perPage := getQueryInt(r, "per_page", 20)

	records, total, err := h.overtimeService.List(userID, month, page, perPage)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	response.SuccessWithMeta(w, http.StatusOK, "Overtime records", records, &response.Meta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

func (h *OvertimeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	id := getPathParam(r, "id")

	record, err := h.overtimeService.GetByID(userID, id)
	if err != nil {
		response.NotFound(w, "Overtime record not found")
		return
	}

	response.OK(w, "Overtime record", record)
}

func (h *OvertimeHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	id := getPathParam(r, "id")

	var input service.OvertimeInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	errors := validator.ValidationErrors{}
	validator.Required(errors, "date", input.Date)
	validator.MinValue(errors, "hours", input.Hours, 0.5)
	validator.MaxValue(errors, "hours", input.Hours, 12)

	if errors.HasErrors() {
		response.ValidationError(w, errors)
		return
	}

	record, err := h.overtimeService.Update(userID, id, input)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, "Overtime record updated", record)
}

func (h *OvertimeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	id := getPathParam(r, "id")

	if err := h.overtimeService.Delete(userID, id); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, "Overtime record deleted", nil)
}

func (h *OvertimeHandler) Calculate(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	hoursStr := r.URL.Query().Get("hours")
	isHolidayStr := r.URL.Query().Get("is_holiday")

	hours := 0.0
	if hoursStr != "" {
		var err error
		hours, err = strToFloat(hoursStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid hours value")
			return
		}
	}

	isHoliday := isHolidayStr == "true" || isHolidayStr == "1"

	calc, err := h.overtimeService.Calculate(userID, hours, isHoliday)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.OK(w, "Overtime calculation", calc)
}

func (h *OvertimeHandler) Disburse(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)

	var input service.DisburseInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	errors := validator.ValidationErrors{}
	validator.Required(errors, "period_start", input.PeriodStart)
	validator.Required(errors, "period_end", input.PeriodEnd)
	validator.Required(errors, "wallet_id", input.WalletID)

	if errors.HasErrors() {
		response.ValidationError(w, errors)
		return
	}

	if err := h.overtimeService.Disburse(userID, input); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, "Period disbursed successfully", nil)
}

func strToFloat(s string) (float64, error) {
	var f float64
	_, err := json.Number(s).Float64()
	if err != nil {
		return 0, err
	}
	f, _ = json.Number(s).Float64()
	return f, nil
}
