package handler

import (
	"encoding/json"
	"net/http"

	"github.com/financeku/backend/internal/service"
	"github.com/financeku/backend/pkg/response"
	"github.com/financeku/backend/pkg/validator"
)

type WalletHandler struct {
	walletService *service.WalletService
}

func NewWalletHandler(walletService *service.WalletService) *WalletHandler {
	return &WalletHandler{walletService: walletService}
}

func (h *WalletHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)

	var input service.WalletInput
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

	wallet, err := h.walletService.Create(userID, input)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.Created(w, "Wallet created", wallet)
}

func (h *WalletHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)

	wallets, err := h.walletService.List(userID)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.OK(w, "Wallets", wallets)
}

func (h *WalletHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	id := getPathParam(r, "id")

	wallet, err := h.walletService.GetByID(userID, id)
	if err != nil {
		response.NotFound(w, "Wallet not found")
		return
	}

	response.OK(w, "Wallet", wallet)
}

func (h *WalletHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	id := getPathParam(r, "id")

	var input service.WalletInput
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

	wallet, err := h.walletService.Update(userID, id, input)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, "Wallet updated", wallet)
}

func (h *WalletHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	id := getPathParam(r, "id")

	if err := h.walletService.Delete(userID, id); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, "Wallet deleted", nil)
}

func (h *WalletHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)

	var input service.TransferInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	errors := validator.ValidationErrors{}
	validator.Required(errors, "from_wallet_id", input.FromWalletID)
	validator.Required(errors, "to_wallet_id", input.ToWalletID)
	validator.MinValue(errors, "amount", input.Amount, 1)

	if errors.HasErrors() {
		response.ValidationError(w, errors)
		return
	}

	if err := h.walletService.Transfer(userID, input); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, "Transfer successful", nil)
}
