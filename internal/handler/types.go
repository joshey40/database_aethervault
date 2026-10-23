package handler

import (
	"net/http"

	"github.com/go-chi/render"
)

type RegisterData struct {
	Username   string `json:"username"`
	Passwd     string `json:"passwd"`
	DeviceName string `json:"currDevice"`
}

func (rd *RegisterData) Bind(r *http.Request) error {
	return nil
}

type LoginData struct {
	RegisterData
}

func (rd *LoginData) Bind(r *http.Request) error {
	return nil
}

type UserCreatedResponse struct {
	AuthToken string `json:"authToken"`
}

func (us *UserCreatedResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type LoginSuccessRespone struct {
	AuthToken string `json:"authToken"`
}

func (us *LoginSuccessRespone) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type ErrResponse struct {
	HTTPStatusCode int    `json:"-"`
	StatusText     string `json:"status"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}
