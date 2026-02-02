package server

import (
	"context"
	"errors"
	"net/http"
	"sellers-accounts-backend/internal/delivery/http/server/middleware"
	"time"

	"sellers-accounts-backend/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HandlersInterface interface {
	RegisterRoutes()
}

type MiddlewareInterface interface {
	Register()
}
type Server struct {
	http     *http.Server
	log      *zap.SugaredLogger
	cfg      *config.Config
	handlers HandlersInterface
	mdl      middleware.Middleware
}

func New(log *zap.SugaredLogger, cfg *config.Config, engine *gin.Engine, handlers HandlersInterface, mdl middleware.Middleware) *Server {
	httpSrv := &http.Server{
		Addr:    cfg.Server.Host + ":" + cfg.Server.Port,
		Handler: engine,
	}

	return &Server{
		http:     httpSrv,
		log:      log,
		cfg:      cfg,
		handlers: handlers,
		mdl:      mdl,
	}
}

func (s *Server) OnStart() error {
	if s.cfg.Server.TLS.Enabled {
		if err := ensureTLSFiles(s.cfg.Server.TLS, s.log); err != nil {
			return err
		}
	}

	s.mdl.Register()
	s.handlers.RegisterRoutes()
	go func() {
		s.log.Infow("server started", "host", s.cfg.Server.Host, "port", s.cfg.Server.Port, "tls", s.cfg.Server.TLS.Enabled)

		var err error
		if s.cfg.Server.TLS.Enabled {
			err = s.http.ListenAndServeTLS(s.cfg.Server.TLS.CertFile, s.cfg.Server.TLS.KeyFile)
		} else {
			err = s.http.ListenAndServe()
		}

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Errorw("failed to serve", "error", err.Error())
		}
	}()
	return nil
}

func (s *Server) OnStop(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.http.Shutdown(shutdownCtx); err != nil {
		return err
	}
	return nil
}
