package handler

import (
	"errors"
	"net/http"

	"sellers-accounts-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

var (
	ErrRequestBody = errors.New("invalid request body")
)

type APIError struct {
	Message string `json:"message"`
}

func writeError(c *gin.Context, err error) {
	status, message := mapRepoError(err)
	if message == "" {
		message = http.StatusText(status)
	}
	c.JSON(status, APIError{Message: message})
}

func mapRepoError(err error) (int, string) {
	switch {
	case err == nil:
		return http.StatusOK, ""
	case errors.Is(err, usecase.ErrNotFoundAccount),
		errors.Is(err, usecase.ErrNotFoundAccountData):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, usecase.ErrNilAccount),
		errors.Is(err, usecase.ErrInvalidAccount),
		errors.Is(err, usecase.ErrInvalidParams):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, usecase.ErrGetAccounts),
		errors.Is(err, ErrRequestBody),
		errors.Is(err, usecase.ErrGetAccount),
		errors.Is(err, usecase.ErrGetAccountData),
		errors.Is(err, usecase.ErrGetTypes),
		errors.Is(err, usecase.ErrCreateAccount),
		errors.Is(err, usecase.ErrUpdateAccount),
		errors.Is(err, usecase.ErrDeleteAccount),
		errors.Is(err, usecase.ErrDeleteAccountData),
		errors.Is(err, usecase.ErrGetDeals),
		errors.Is(err, usecase.ErrInternal):
		return http.StatusInternalServerError, err.Error()
	default:
		return http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)
	}
}
