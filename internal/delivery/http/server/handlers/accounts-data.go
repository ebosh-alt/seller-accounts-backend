package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"sellers-accounts-backend/internal/entities"
	"strconv"
)

// AccountData
// @Summary Get account data by ID
// @Tags accounts
// @Param id path string true "Data ID"
// @Success 200 {object} entities.ResponseGetAccountData
// @Failure 400 {object} APIError
// @Failure 404 {object} APIError
// @Failure 500 {object} APIError
// @Router /accounts/data/{id} [get]
func (h *Handler) AccountData(c *gin.Context) {
	param := c.Param("id")
	if param == "" {
		h.log.Warnw("failed to bind get account data request")
		writeBadRequest(c, "invalid request body")
		return
	}
	dataID, err := strconv.Atoi(param)
	if err != nil {
		fmt.Println(err)
	}
	req := entities.RequestGetAccountData{ID: dataID}

	resp, err := h.uc.AccountData(c.Request.Context(), &req)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(200, resp)
}

// DeleteAccountData
// @Summary delete account data
// @Tags data
// @Accept json
// @Produce json
// @Param id path string true "Data ID"
// @Success 200 {object} entities.ResponseDeleteAccountData
// @Failure 400 {object} APIError
// @Failure 500 {object} APIError
// @Router /accounts/data/{id} [delete]
func (h *Handler) DeleteAccountData(c *gin.Context) {
	param := c.Param("id")
	if param == "" {
		h.log.Warnw("failed to bind get account data request")
		writeBadRequest(c, "invalid request body")
		return
	}
	dataID, err := strconv.Atoi(param)
	if err != nil {
		h.log.Warnw("failed to bind get account data request")
		writeBadRequest(c, "invalid request body")
		return
	}
	req := entities.RequestDeleteAccountData{ID: dataID}
	resp, err := h.uc.DeleteAccountData(c.Request.Context(), &req)
	if err != nil {
		writeError(c, err)
		return
	}
	h.log.Infow("deleted account data", "resp", resp)
	c.JSON(200, resp)
}
