package handler

import (
	"net/http"

	"sellers-accounts-backend/internal/usecase"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

type InterfaceHandler interface {
	AllAccounts(c *gin.Context)
	AcceptableTypesAccounts(c *gin.Context)
	Account(c *gin.Context)
	CreateAccounts(c *gin.Context)
	// UpdateAccount CreateAccounts(c *gin.Context)
	UpdateAccount(c *gin.Context)
	AccountData(c *gin.Context)
	DeleteAccountData(c *gin.Context)
	Deals(c *gin.Context)
	RegisterRoutes()
}

type Handler struct {
	engine *gin.Engine
	log    *zap.SugaredLogger
	uc     usecase.InterfaceUsecase
}

func New(log *zap.SugaredLogger, engine *gin.Engine, uc usecase.InterfaceUsecase) InterfaceHandler {
	return &Handler{log: log, engine: engine, uc: uc}
}

func (h *Handler) RegisterRoutes() {
	h.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := h.engine.Group("/api")
	{
		accountsApi := api.Group("/accounts")
		{
			accountsApi.GET("/all", h.AllAccounts)
			accountsApi.GET("/:id", h.Account)
			accountsApi.GET("/:id/data")

			accountsApi.POST("/", h.CreateAccounts)
			accountsApi.PUT("/", h.UpdateAccount)
			//accountsApi.POST("/deactivate-view-by-name", h.DeactivateAccountsByName)
		}

		typesApi := api.Group("/acceptable-types")
		{
			typesApi.GET("/", h.AcceptableTypesAccounts)
		}

		dataApi := accountsApi.Group("/data")
		{
			dataApi.GET("/:id", h.AccountData)

			dataApi.DELETE("/:id", h.DeleteAccountData)

			dataApi.PUT("/update")
		}
		dealsApi := api.Group("/deals")
		{
			dealsApi.GET("/", h.Deals)
		}
	}

}

func writeBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, APIError{Message: message})
}
