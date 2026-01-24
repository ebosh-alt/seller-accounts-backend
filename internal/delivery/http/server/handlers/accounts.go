package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"sellers-accounts-backend/internal/entities"
)

// AllAccounts
// @Summary List accounts
// @Description Returns a paginated list of accounts.
// @Tags accounts
// @Param limit query int false "Page size"
// @Param page query int false "Page number (1-based)"
// @Success 200 {object} entities.ResponseAllAccounts
// @Failure 400 {object} APIError
// @Failure 500 {object} APIError
// @Router /accounts/all [get]
func (h *Handler) AllAccounts(c *gin.Context) {
	var req entities.RequestAllAccounts
	if err := c.ShouldBindQuery(&req); err != nil {
		h.log.Warnw("failed to bind list accounts request", zap.Error(err))
		writeBadRequest(c, "invalid request params")
		return
	}

	response, err := h.uc.AllAccounts(c.Request.Context(), &req)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// Account
// @Summary Get account by ID
// @Tags accounts
// @Param id path string true "Account ID"
// @Success 200 {object} entities.ResponseAccountByID
// @Failure 400 {object} APIError
// @Failure 404 {object} APIError
// @Failure 500 {object} APIError
// @Router /accounts/{id} [get]
func (h *Handler) Account(c *gin.Context) {
	accID := c.Param("id")
	if accID == "" {
		h.log.Warnw("failed to bind get account request")
		writeBadRequest(c, "invalid request body")
		return
	}
	req := entities.RequestAccountByID{ID: accID}
	h.log.Infof("id: %s", req.ID)
	response, err := h.uc.Account(c.Request.Context(), &req)
	if err != nil {
		h.log.Warnw("error", err)
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateAccounts
// @Summary Create account
// @Tags accounts
// @Accept json
// @Produce json
// @Param request body entities.RequestCreateAccounts true "Create accounts request"
// @Success 201 {object} entities.ResponseCreateAccounts
// @Failure 400 {object} APIError
// @Failure 500 {object} APIError
// @Router /accounts/ [post]
func (h *Handler) CreateAccounts(c *gin.Context) {
	var req entities.RequestCreateAccounts
	if err := c.ShouldBindJSON(&req); err != nil {
		writeBadRequest(c, "invalid request body")
		return
	}
	response, err := h.uc.CreateAccounts(c.Request.Context(), &req)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// UpdateAccount
// @Summary Update account
// @Tags accounts
// @Accept json
// @Produce json
// @Param request body entities.RequestUpdateAccount true "Update account request"
// @Success 200 {object} entities.ResponseUpdateAccount
// @Failure 400 {object} APIError
// @Failure 500 {object} APIError
// @Router /accounts [put]
func (h *Handler) UpdateAccount(c *gin.Context) {
	var req entities.RequestUpdateAccount
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Infoln(err)
		writeBadRequest(c, "invalid request body")
		return
	}
	response, err := h.uc.UpdateAccount(c.Request.Context(), &req)
	if err != nil {
		h.log.Warnw("err", err)
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
