package handler

import (
	"encoding/json"
	"net/http"

	"github.com/financeku/backend/internal/repository"
	"github.com/financeku/backend/internal/service"
	"github.com/financeku/backend/pkg/response"
	"github.com/financeku/backend/pkg/validator"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
}

func NewTransactionHandler(transactionService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)

	var input service.TransactionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	errors := validator.ValidationErrors{}
	validator.Required(errors, "wallet_id", input.WalletID)
	validator.Required(errors, "type", input.Type)
	validator.InList(errors, "type", input.Type, []string{"income", "expense"})
	validator.MinValue(errors, "amount", input.Amount, 1)
	validator.Required(errors, "date", input.Date)

	if errors.HasErrors() {
		response.ValidationError(w, errors)
		return
	}

	tx, err := h.transactionService.Create(userID, input)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Created(w, "Transaction created", tx)
}

func (h *TransactionHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)

	filter := repository.TransactionFilter{
		UserID:     userID,
		WalletID:   getQueryString(r, "wallet_id", ""),
		CategoryID: getQueryString(r, "category_id", ""),
		Type:       getQueryString(r, "type", ""),
		DateFrom:   getQueryString(r, "date_from", ""),
		DateTo:     getQueryString(r, "date_to", ""),
		Page:       getQueryInt(r, "page", 1),
		PerPage:    getQueryInt(r, "per_page", 20),
	}

	transactions, total, err := h.transactionService.List(userID, filter)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	totalPages := int(total) / filter.PerPage
	if int(total)%filter.PerPage > 0 {
		totalPages++
	}

	response.SuccessWithMeta(w, http.StatusOK, "Transactions", transactions, &response.Meta{
		Page:       filter.Page,
		PerPage:    filter.PerPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

func (h *TransactionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	id := getPathParam(r, "id")

	tx, err := h.transactionService.GetByID(userID, id)
	if err != nil {
		response.NotFound(w, "Transaction not found")
		return
	}

	response.OK(w, "Transaction", tx)
}

func (h *TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	id := getPathParam(r, "id")

	var input service.TransactionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	errors := validator.ValidationErrors{}
	validator.Required(errors, "wallet_id", input.WalletID)
	validator.Required(errors, "type", input.Type)
	validator.InList(errors, "type", input.Type, []string{"income", "expense"})
	validator.MinValue(errors, "amount", input.Amount, 1)
	validator.Required(errors, "date", input.Date)

	if errors.HasErrors() {
		response.ValidationError(w, errors)
		return
	}

	tx, err := h.transactionService.Update(userID, id, input)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, "Transaction updated", tx)
}

func (h *TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	id := getPathParam(r, "id")

	if err := h.transactionService.Delete(userID, id); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, "Transaction deleted", nil)
}
