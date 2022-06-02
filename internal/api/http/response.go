package http

import (
	"context"
	"encoding/json"
	"net/http"

	"webChat/internal/ctxkey"
	ierr "webChat/internal/errors"
)

const (
	HeaderContentType   = "Content-Type"
	MimeApplicationJSON = "application/json"
	MimeTextPlain       = "text/plain"
)

type ResponseManager struct{}

func NewResponseManager() *ResponseManager {
	return &ResponseManager{}
}

func (rm *ResponseManager) JSON(ctx context.Context, w http.ResponseWriter, code int, data interface{}) {
	if data == nil {
		w.WriteHeader(code)
		return
	}

	w.Header().Set(HeaderContentType, MimeApplicationJSON)
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(&data); err != nil {
		ctxkey.GetLogger(ctx).Error("failed to respond with JSON", err)
	}
}

type ApiError struct {
	Error string `json:"error"`
}

func (rm *ResponseManager) Error(ctx context.Context, w http.ResponseWriter, err error) {
	e := ierr.Get(err)
	code := e.Code()

	if code >= http.StatusInternalServerError {
		ctxkey.GetLogger(ctx).Error("got an api error", err)
		w.WriteHeader(code)
		return
	}

	rm.JSON(ctx, w, code, ApiError{Error: err.Error()})
}

func (rm *ResponseManager) OK(ctx context.Context, w http.ResponseWriter, data interface{}) {
	rm.JSON(ctx, w, http.StatusOK, data)
}

func (rm *ResponseManager) Created(ctx context.Context, w http.ResponseWriter, data interface{}) {
	rm.JSON(ctx, w, http.StatusCreated, data)
}

func (rm *ResponseManager) Write(ctx context.Context, w http.ResponseWriter, code int, mime string, data []byte) {
	if data == nil {
		w.WriteHeader(code)
		return
	}

	w.Header().Set(HeaderContentType, mime)
	w.WriteHeader(code)

	if _, err := w.Write(data); err != nil {
		ctxkey.GetLogger(ctx).Error("failed to write response", err)
	}
}
