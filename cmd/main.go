package main

import (
	"TaskTracker/internal/config"
	"TaskTracker/internal/pkg/db"
	"TaskTracker/internal/pkg/middleware"
	"TaskTracker/internal/pkg/server"
	"TaskTracker/internal/pkg/services/auth"
	authHandler "TaskTracker/internal/pkg/services/auth/delivery/http"
	authRepo "TaskTracker/internal/pkg/services/auth/repo"
	authUsecase "TaskTracker/internal/pkg/services/auth/usecase"
	"TaskTracker/internal/pkg/services/board"
	boardHandler "TaskTracker/internal/pkg/services/board/delivery/http"
	boardRepo "TaskTracker/internal/pkg/services/board/repo"
	boardUsecase "TaskTracker/internal/pkg/services/board/usecase"
	"TaskTracker/internal/pkg/services/user"
	userHandler "TaskTracker/internal/pkg/services/user/delivery/http"
	userRepo "TaskTracker/internal/pkg/services/user/repo"
	userUsecase "TaskTracker/internal/pkg/services/user/usecase"
	"TaskTracker/internal/pkg/tokenizer"
	"TaskTracker/pkg/logger"
	"context"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := fx.New(
		fx.Provide(
			logger.SetupLogger,
			server.NewRouter,

			config.MustLoad,

			db.NewClickHouse,

			tokenizer.New,
			middleware.NewAuthMiddleware,

			authHandler.New,
			fx.Annotate(authUsecase.New, fx.As(new(auth.Usecase))),
			fx.Annotate(authRepo.New, fx.As(new(auth.Repository))),

			userHandler.New,
			fx.Annotate(userUsecase.New, fx.As(new(user.Usecase))),
			fx.Annotate(userRepo.New, fx.As(new(user.Repository))),

			boardHandler.New,
			fx.Annotate(boardUsecase.New, fx.As(new(board.Usecase))),
			fx.Annotate(boardRepo.New, fx.As(new(board.Repository))),
		),

		fx.WithLogger(func(logger *slog.Logger) fxevent.Logger {
			return &fxevent.SlogLogger{Logger: logger}
		}),

		fx.Invoke(
			server.RunServer,
			//migrations.RunMigrations,
		),
	)

	ctx := context.Background()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	if err := app.Start(ctx); err != nil {
		panic(err)
	}

	<-stop
	app.Stop(ctx)
}
