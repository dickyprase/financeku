package handler

import (
	"encoding/json"
	"net/http"

	"github.com/financeku/backend/internal/service"
	"github.com/financeku/backend/pkg/response"
	"github.com/financeku/backend/pkg/validator"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)

	var input service.CategoryInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	errors := validator.ValidationErrors{}
	validator.Required(errors, "name", input.Name)
	validator.Required(errors, "type", input.Type)
	validator.InList(errors, "type", input.Type, []string{"income", "expense"})

	if errors.HasErrors() {
		response.ValidationError(w, errors)
		return
	}

	cat, err := h.categoryService.Create(userID, input)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.Created(w, "Category created", cat)
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	catType := getQueryString(r, "type", "")

	categories, err := h.categoryService.List(userID, catType)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.OK(w, "Categories", categories)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	id := getPathParam(r, "id")

	var input service.CategoryInput
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

	cat, err := h.categoryService.Update(userID, id, input)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, "Category updated", cat)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	id := getPathParam(r, "id")

	if err := h.categoryService.Delete(userID, id); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, "Category deleted", nil)
}
