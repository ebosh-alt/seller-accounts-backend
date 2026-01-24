package middleware

import (
	"sellers-accounts-backend/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Middleware struct {
	cfg    *config.Config
	log    *zap.SugaredLogger
	engine *gin.Engine
}

func New(cfg *config.Config, log *zap.SugaredLogger, engine *gin.Engine) *Middleware {
	return &Middleware{
		cfg:    cfg,
		log:    log,
		engine: engine,
	}
}

func (m *Middleware) Register() {
	m.engine.Use(m.CORSMiddleware())
	m.log.Info("register middleware")
}
