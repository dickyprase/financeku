package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type ErrorResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func Success(w http.ResponseWriter, status int, message string, data interface{}) {
	JSON(w, status, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessWithMeta(w http.ResponseWriter, status int, message string, data interface{}, meta *Meta) {
	JSON(w, status, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, ErrorResponse{
		Success: false,
		Message: message,
	})
}

func ValidationError(w http.ResponseWriter, errors map[string]string) {
	JSON(w, http.StatusUnprocessableEntity, ErrorResponse{
		Success: false,
		Message: "Validation failed",
		Errors:  errors,
	})
}

func Created(w http.ResponseWriter, message string, data interface{}) {
	Success(w, http.StatusCreated, message, data)
}

func OK(w http.ResponseWriter, message string, data interface{}) {
	Success(w, http.StatusOK, message, data)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func Unauthorized(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Unauthorized"
	}
	Error(w, http.StatusUnauthorized, message)
}

func Forbidden(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Forbidden"
	}
	Error(w, http.StatusForbidden, message)
}

func NotFound(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Resource not found"
	}
	Error(w, http.StatusNotFound, message)
}

func InternalError(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Internal server error"
	}
	Error(w, http.StatusInternalServerError, message)
}
