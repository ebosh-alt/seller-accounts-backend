package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"sellers-accounts-backend/internal/entities"
)

// Deals
// @Summary List deals
// @Description Returns a paginated list of deals.
// @Tags deals
// @Param limit query int false "Page size"
// @Param page query int false "Page number (1-based)"
// @Success 200 {object} entities.ResponseGetDeals
// @Failure 400 {object} APIError
// @Failure 500 {object} APIError
// @Router /deals [get]
func (h *Handler) Deals(c *gin.Context) {
	var req entities.RequestGetDeals
	if err := c.ShouldBindQuery(&req); err != nil {
		h.log.Warnw("failed to bind list deals request", zap.Error(err))
		writeBadRequest(c, "invalid request params")
		return
	}

	response, err := h.uc.Deals(c.Request.Context(), &req)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
