package http

import (
	"TaskTracker/internal/models"
	"TaskTracker/internal/pkg/constans"
	"TaskTracker/internal/pkg/services/auth"
	"TaskTracker/pkg/reader"
	"TaskTracker/pkg/responser"
	"go.uber.org/fx"
	"log/slog"
	"net/http"
)

type Params struct {
	fx.In

	Usecase auth.Usecase
	Logger  *slog.Logger
}

type Handler struct {
	useCase auth.Usecase
	log     *slog.Logger
}

func New(params Params) *Handler {
	return &Handler{
		useCase: params.Usecase,
		log:     params.Logger,
	}
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var userData *models.SignUpRequest

	if err := reader.ReadRequestData(r, &userData); err != nil {
		h.log.Error("read request data: ", err.Error())
		responser.Send400(w, "error to read request data")
		return
	}

	token, err := h.useCase.CreateUser(r.Context(), userData)
	if err != nil {
		h.log.Error("error to create employee: ", err.Error())
		responser.Send400(w, "error to create employee")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     constans.CookieName,
		Value:    token.Token,
		Path:     "/",
		Expires:  token.Exp,
		HttpOnly: true,
	})

	h.log.Debug("success to create employee: " + token.Token)
	responser.Send201(w, token)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var userData *models.LoginRequest
	if err := reader.ReadRequestData(r, &userData); err != nil {
		h.log.Error("read request data: " + err.Error())
		responser.Send400(w, "error to read request data")
		return
	}

	token, err := h.useCase.Login(r.Context(), userData)
	if err != nil {
		h.log.Error("error to login: " + err.Error())
		responser.Send400(w, "error to read request data")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     constans.CookieName,
		Value:    token.Token,
		Path:     "/",
		Expires:  token.Exp,
		HttpOnly: true,
	})

	responser.Send200(w, token)
}
