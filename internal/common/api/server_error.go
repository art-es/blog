//go:generate mockgen -source=server_error.go -destination=mock/server_error.go -package=mock
package api

import (
	"encoding/json"
	"net/http"

	"github.com/art-es/blog/internal/common/log"
)

var _ ServerErrorHandler = (*serverErrorHandler)(nil)

type ServerErrorHandler interface {
	Handle(endpoint string, w http.ResponseWriter, err error)
}

const responseMessage = "sorry, please, try again later"

type response struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

type serverErrorHandler struct {
	logger log.Logger
}

func NewServerErrorHandler(logger log.Logger) *serverErrorHandler {
	return &serverErrorHandler{
		logger: logger,
	}
}

func (h *serverErrorHandler) Handle(endpoint string, w http.ResponseWriter, err error) {
	h.logger.Error("API server error", log.String("endpointName", endpoint), log.Error(err))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(&response{
		OK:      false,
		Message: responseMessage,
	})
}
