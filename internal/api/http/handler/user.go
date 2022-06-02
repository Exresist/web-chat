package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) getUserByID(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("user_id")
	id, err := strconv.Atoi(param)
	if err != nil {
		s.respond.Error(ctx, c.Writer, err)
		return
	}

	user, err := s.userStore.GetByID(ctx, id)
	if err != nil {
		s.respond.Error(ctx, c.Writer, err)
		return
	}

	s.respond.JSON(ctx, c.Writer, http.StatusOK, user)

}
