package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/config"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/graphql"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/middlewares"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/repository"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/service"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	inmemoryStorage = "inmemory"
	postgresStorage = "postgres"
)

type App struct {
	log    *slog.Logger
	config *config.Config
	server *http.Server
}

func New(log *slog.Logger, config *config.Config) *App {
	if config == nil {
		panic("config is nil")
	}

	var repo repository.Repository
	if config.StorageType == inmemoryStorage {
		repo = repository.NewInMemoryRepository()
	} else if config.StorageType == postgresStorage {

	} else {
		panic("unsupported storage type")
	}

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", config.Port),
	}

	postService := service.NewPostService(repo)
	userService := service.NewUserService(repo)
	commentService := service.NewCommentService(repo)

	app := &App{
		log:    log,
		config: config,
		server: server,
	}
	app.SetHandlers(repo, postService, userService, commentService)
	return app
}

func (a *App) SetHandlers(repo repository.Repository,
	postService *service.PostService,
	userService *service.UserService,
	commentService *service.CommentService) {

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	rootHandler := middlewares.DataloaderMiddleware(repo,
		handler.NewDefaultServer(
			graphql.NewExecutableSchema(graphql.NewRootResolver(postService, userService, commentService)),
		),
	)
	http.Handle("/query", middlewares.AuthMiddleware(rootHandler))
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "app.Run"

	log := a.log.With(
		slog.String("op", op),
	)

	stop := make(chan os.Signal, 1)
	errChan := make(chan error, 1)

	go func() {
		log.Info("Server starting", slog.String("address", a.server.Addr))
		if err := http.ListenAndServe(a.server.Addr, nil); !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
			return
		}
		log.Info("Server stopped")
	}()

	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <-errChan:
		log.Error(err.Error())
		return err
	case <-stop:
		log.Info("Stop signal received, shutting down gracefully")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("Graceful shutdown complete")
	return nil
}
