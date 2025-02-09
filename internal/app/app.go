package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/config"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/graphql"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/middlewares"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/repository"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/service"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/utils"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	if config.StorageType == utils.InmemoryStorage {
		repo = repository.NewInMemoryRepository()
	} else if config.StorageType == utils.PostgresStorage {
		// TODO: postgres
	} else {
		panic("unsupported storage type")
	}

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", config.Port),
	}

	postService := service.NewPostService(repo)
	userService := service.NewUserService(repo)
	commentService := service.NewCommentService(repo)
	pubSub := service.NewPubSub()

	app := &App{
		log:    log,
		config: config,
		server: server,
	}
	app.SetupHandlers(repo, postService, userService, commentService, pubSub)
	return app
}

func (a *App) SetupHandlers(repo repository.Repository,
	postService *service.PostService,
	userService *service.UserService,
	commentService *service.CommentService,
	pubSub *service.PubSub) {

	if a.config.Env != utils.EnvProd {
		http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	}

	defaultServer := handler.NewDefaultServer(
		graphql.NewExecutableSchema(graphql.NewRootResolver(postService, userService, commentService, pubSub)),
	)
	defaultServer.Use(extension.FixedComplexityLimit(100))
	rootHandler := middlewares.DataloaderMiddleware(repo, defaultServer)

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
		if a.config.Env != utils.EnvProd {
			log.Info("Playground is running", slog.String("address", a.server.Addr))
		}
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
