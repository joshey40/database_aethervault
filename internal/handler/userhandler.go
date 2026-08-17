package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/joshey40/database_aethervault/internal/logger"
	"github.com/joshey40/database_aethervault/internal/services"
	"go.uber.org/zap"
)

type UserHandler struct {
	sv *services.UserService
}

func NewUserHandler(sv *services.UserService) *UserHandler {
	return &UserHandler{
		sv: sv,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	cont := new(strings.Builder)
	n, err := io.Copy(cont, r.Body)
	if err != nil {
		logger.L().Error("Reading Body of CreateUser api call failed", zap.Error(err))
	}
	logger.L().Info("Reading of Body successful", zap.String("Body content", cont.String()), zap.Int("Num Bytes read", int(n)))
}
