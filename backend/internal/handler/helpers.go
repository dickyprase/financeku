package handler

import (
	"net/http"
	"strconv"

	"github.com/financeku/backend/internal/middleware"
)

func getUserIDFromContext(r *http.Request) string {
	return middleware.GetUserID(r)
}

func getPathParam(r *http.Request, name string) string {
	return r.PathValue(name)
}

func getQueryInt(r *http.Request, key string, defaultVal int) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return i
}

func getQueryString(r *http.Request, key string, defaultVal string) string {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultVal
	}
	return val
}
