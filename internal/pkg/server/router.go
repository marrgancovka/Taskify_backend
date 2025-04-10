package server

import (
	"TaskTracker/internal/pkg/middleware"
	authHandler "TaskTracker/internal/pkg/services/auth/delivery/http"
	boardHandler "TaskTracker/internal/pkg/services/board/delivery/http"
	userHandler "TaskTracker/internal/pkg/services/user/delivery/http"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"log/slog"
	"net/http"
)

type RouterParams struct {
	fx.In

	AuthMiddleware *middleware.AuthMiddleware

	AuthHandler  *authHandler.Handler
	UserHandler  *userHandler.Handler
	BoardHandler *boardHandler.Handler
	Logger       *slog.Logger
}

type Router struct {
	handler *mux.Router
}

func NewRouter(p RouterParams) *Router {
	api := mux.NewRouter().PathPrefix("/api").Subrouter()
	api.Use(middleware.CORSMiddleware)

	v1 := api.PathPrefix("/v1").Subrouter()

	auth := v1.PathPrefix("/auth").Subrouter()

	auth.HandleFunc("/signup", p.AuthHandler.SignUp).Methods(http.MethodPost, http.MethodOptions)
	auth.HandleFunc("/login", p.AuthHandler.Login).Methods(http.MethodPost, http.MethodOptions)

	users := v1.PathPrefix("/user").Subrouter()
	users.Use(p.AuthMiddleware.JwtMiddleware)

	users.HandleFunc("/me", p.UserHandler.GetUserByID).Methods(http.MethodGet, http.MethodOptions)

	boards := v1.PathPrefix("/board").Subrouter()
	boards.Use(p.AuthMiddleware.JwtMiddleware)

	boards.HandleFunc("", p.BoardHandler.GetUserListBoard).Methods(http.MethodGet, http.MethodOptions)
	boards.HandleFunc("", p.BoardHandler.CreateBoard).Methods(http.MethodPost, http.MethodOptions)
	boards.HandleFunc("/{boardID}/tasks", p.BoardHandler.GetBoardTasks).Methods(http.MethodGet, http.MethodOptions)
	boards.HandleFunc("/favourite/{boardID}", p.BoardHandler.SetFavouriteBoard).Methods(http.MethodPut, http.MethodOptions)
	boards.HandleFunc("/nofavourite/{boardID}", p.BoardHandler.SetNoFavouriteBoard).Methods(http.MethodPut, http.MethodOptions)
	boards.HandleFunc("/{boardID}/member", p.BoardHandler.AddMember).Methods(http.MethodPost, http.MethodOptions)

	router := &Router{
		handler: api,
	}

	p.Logger.Info("registered router")

	return router
}
