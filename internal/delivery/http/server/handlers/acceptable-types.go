package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// AcceptableTypesAccounts
// @Summary List acceptable account types
// @Tags accounts
// @Success 200 {object} entities.ResponseAcceptableTypeAccounts
// @Failure 500 {object} APIError
// @Router /acceptable-types [get]
func (h *Handler) AcceptableTypesAccounts(c *gin.Context) {
	response, err := h.uc.AcceptableTypesAccounts(c.Request.Context())
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
