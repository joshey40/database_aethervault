package handler

import (
	"net/http"

	"github.com/go-chi/render"
)

type RegisterData struct {
	Username string `json:"username"`
	Passwd   string `json:"passwd"`
}

func (rd *RegisterData) Bind(r *http.Request) error {
	return nil
}

type UserCreatedResponse struct {
}

func (us *UserCreatedResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}
