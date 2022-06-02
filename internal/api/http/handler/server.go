package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/run"

	api "webChat/internal/api/http"
	"webChat/internal/config"
	"webChat/internal/db"

	log "go.uber.org/zap"
)

const (
	headerXRequestID      = "X-Request-Id"
	constAuthHeader       = "Authorization"
	constBearerAuthPrefix = "Bearer"
	constContentType      = "Content-Type"
)

type Server struct {
	*http.Server
	respond *api.ResponseManager
	cfg     *config.Config
	logger  *log.SugaredLogger

	userStore db.UserStore
}

func NewServer(cfg *config.Config, logger *log.SugaredLogger,
	userStore db.UserStore) *Server {
	srv := &Server{
		Server: &http.Server{
			Addr:         cfg.API.Address,
			ReadTimeout:  cfg.API.ReadTimeout,
			WriteTimeout: cfg.API.WriteTimeout,
		},
		cfg:       cfg,
		logger:    logger,
		userStore: userStore,
	}

	r := gin.Default()
	r.Use(ginLogger(logger))
	r.Use(gin.Recovery())

	r.POST("/auth", srv.auth)
	r.POST("/register", srv.register)

	r.Group("/user").
		Use(srv.authMiddleware).
		GET("/{user_id}", srv.getUserByID)

	srv.Handler = r
	return srv
}

func (s *Server) Run(g *run.Group) {
	g.Add(func() error {
		s.logger.Info("[http-server] started")
		return s.ListenAndServe()
	}, func(err error) {
		s.logger.Error("[http-server] terminated", err)

		ctx, cancel := context.WithTimeout(context.Background(), s.cfg.API.ShutdownTimeout)
		defer cancel()

		s.logger.Error("[http-server] stopped", s.Shutdown(ctx))
	})
}
