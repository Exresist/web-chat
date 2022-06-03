package handler

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "go.uber.org/zap"

	service "webChat/internal/controller/jwt"
	"webChat/internal/ctxkey"
	ierr "webChat/internal/errors"
)

// GinLogger receives the default log of the gin framework
func ginLogger(logger *log.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		reqID := c.Request.Header.Get(headerXRequestID)
		if reqID == "" {
			reqID = uuid.New().String()
		}
		c.Writer.Header().Set(headerXRequestID, reqID)

		ctx = ctxkey.PutRequestID(ctx, reqID)
		ctx = ctxkey.PutLogger(ctx, logger.With("request_id", reqID))
	}
}

func (s *Server) authMiddleware(c *gin.Context) {
	ctx := c.Request.Context()

	authorization := c.Request.Header.Get(constAuthHeader)
	if authorization == "" {
		s.respond.Error(ctx, c.Writer, ierr.ErrUnauthorized)
		return
	}

	parts := strings.Split(strings.TrimSpace(authorization), " ")
	if len(parts) < 2 || parts[0] != constBearerAuthPrefix {
		s.respond.Error(ctx, c.Writer, ierr.ErrUnauthorized)
		return
	}
	if err := service.ParseToken(parts[len(parts)-1], []byte(s.cfg.SecretKey)); err != nil {
		s.respond.Error(ctx, c.Writer, ierr.ErrUnauthorized)
		return
	}

	c.Next()
}
