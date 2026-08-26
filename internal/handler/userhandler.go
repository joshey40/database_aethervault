package handler

import (
	"net/http"

	"github.com/go-chi/render"
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

	newUser := &RegisterData{}
	err := render.Bind(r, newUser)
	if err != nil {
		logger.L().Error("Reading Body of CreateUser api call failed", zap.Error(err))
		render.Render(w, r, ErrInvalidRequest(err))
	}
	logger.L().Info("Reading of Body successful", zap.String("username", newUser.Username), zap.String("passwd", newUser.Passwd))
	// create user here
	_, err = h.sv.CreateUser(r.Context(), newUser.Username, newUser.Passwd, 0)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, UserCreated())
}

func UserCreated() *UserCreatedResponse {
	return &UserCreatedResponse{}
}

func ErrInvalidRequest(err error) *ErrResponse {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request",
		ErrorText:      err.Error(),
	}
}
