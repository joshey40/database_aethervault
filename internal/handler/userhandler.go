package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
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
		zap.L().Error("Reading Body of CreateUser api call failed", zap.Error(err))
		render.Render(w, r, MakeErrResponse(400, services.ErrInvalidRequest.Error()))
		return
	}
	zap.L().Debug("Reading of Body successful", zap.String("username", newUser.Username), zap.String("passwd", newUser.Passwd))

	// create new user
	token, err := h.sv.CreateUser(r.Context(), newUser.Username, newUser.Passwd, 0, newUser.DeviceName)
	if err != nil {
		zap.L().Warn("Error creating new User", zap.Error(err))
		if errors.Is(err, services.ErrUserAlreadyExists) {
			render.Render(w, r, MakeErrResponse(409, err.Error()))
		} else {
			render.Render(w, r, MakeErrResponse(500, services.ErrInternalError.Error()))
		}
		return
	}
	response := &UserCreatedResponse{
		AuthToken: token,
	}
	render.Status(r, http.StatusCreated)
	render.Render(w, r, response)
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	loginData := &LoginData{}
	err := render.Bind(r, loginData)
	if err != nil {
		zap.L().Error("Reading Body of Login api call failed", zap.Error(err))
		render.Render(w, r, MakeErrResponse(400, services.ErrInvalidRequest.Error()))
		return
	}

	zap.L().Debug("Reading of Body successful", zap.String("username", loginData.Username), zap.String("passwd", loginData.Passwd))

	token, err := h.sv.LoginUser(r.Context(), loginData.Username, loginData.Passwd, loginData.DeviceName)
	if err != nil {
		if errors.Is(err, services.ErrInvalidLoginData) {
			render.Render(w, r, MakeErrResponse(401, err.Error()))
		} else {
			render.Render(w, r, MakeErrResponse(500, services.ErrInternalError.Error()))
		}
		return
	}
	response := &LoginSuccessRespone{
		AuthToken: token,
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, response)
}

func MakeErrResponse(status int, text string) *ErrResponse {
	return &ErrResponse{
		HTTPStatusCode: status,
		StatusText:     text,
	}
}
