package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"webChat/internal/config"
	service "webChat/internal/controller/jwt"
	ierr "webChat/internal/errors"
	"webChat/internal/model"
)

type (
	AuthRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	AuthResponse struct {
		Token string `json:"token"`
	}
	RegisterRequest struct {
		Username  string         `json:"username"`
		Password  string         `json:"password"`
		Photo     multipart.File `json:"photo,omitempty"`
		FirstName string         `json:"first_name,omitempty"`
		LastName  string         `json:"last_name,omitempty"`
		Email     string         `json:"email,omitempty"`
	}
)

func (s *Server) auth(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req AuthRequest
	)

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		s.respond.Error(ctx, c.Writer, err)
		return
	}

	user, err := s.userStore.GetByUsername(ctx, req.Username)
	if err != nil {
		s.respond.Error(ctx, c.Writer, err)
		return
	}
	if user == nil {
		s.respond.Error(ctx, c.Writer, ierr.ErrUserNotFound)
	}
	err = user.ComparePassword(req.Password)
	if err != nil {
		s.respond.Error(ctx, c.Writer, err)
	}

	token, err := service.GenerateToken(s.cfg.SecretKey)
	if err != nil {
		return
	}

	s.respond.OK(ctx, c.Writer, AuthResponse{Token: token})
}

func (s *Server) register(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		req RegisterRequest
	)
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		s.respond.Error(ctx, c.Writer, err)
		return
	}
	if err := c.Request.ParseMultipartForm(0); err != nil {
		s.respond.Error(ctx, c.Writer, err)
		return
	}

	photo, header, err := c.Request.FormFile("photo")
	if err != nil {
		s.respond.Error(ctx, c.Writer, err)
		return
	}
	if err := checkFile(header, s.cfg.Files); err != nil {
		s.respond.Error(ctx, c.Writer, err)
		return
	}
	req.Photo = photo
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, photo); err != nil {
		s.respond.Error(ctx, c.Writer, err)
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), 0)
	if err != nil {
		s.respond.Error(ctx, c.Writer, err)
		return
	}

	user := &model.User{
		Username:       req.Username,
		HashedPassword: string(hashedPass),
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Photo:          buf.Bytes(),
		Email:          req.Email,
	}
	if err := s.userStore.Insert(ctx, user); err != nil {
		s.respond.Error(ctx, c.Writer, err)
		return
	}
	s.respond.OK(ctx, c.Writer, nil)
}

// checkFile checks file content-type and size.
func checkFile(fileHeader *multipart.FileHeader, cfg config.Files) error {
	switch contentType, ok := fileHeader.Header[constContentType]; {
	case !ok:
		return ierr.New("no Content-Type header")
	case fileHeader.Size == 0:
		return ierr.New("empty file")
	case fileHeader.Size > cfg.MaxFileSize:
		return fmt.Errorf("file is too big (max %d)", cfg.MaxFileSize)
	case len(contentType) > 0:
		for _, cfgContType := range cfg.AllowedFileTypes {
			if strings.EqualFold(cfgContType, contentType[0]) {
				return nil
			}
		}
		return ierr.New("wrong file type")
	default:
		return ierr.New("can't access file type")
	}
}
