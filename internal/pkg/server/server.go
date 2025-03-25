package server

import (
	"go.uber.org/fx"
	"net/http"
)

type ServerParams struct {
	fx.In

	Config Config
	Router *Router
}

func RunServer(p ServerParams) {

	srv := &http.Server{
		Addr:              p.Config.Address,
		Handler:           p.Router.handler,
		ReadHeaderTimeout: p.Config.ReadHeaderTimeout,
		IdleTimeout:       p.Config.IdleTimeout,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
}
