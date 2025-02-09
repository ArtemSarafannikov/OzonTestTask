package graphqlapp

import (
	"fmt"
	"log/slog"
	"net/http"
)

func (a *App) Run() error {
	const op = "graphqlapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", a.port), nil); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("connect to http://localhost:%s/ for GraphQL playground", port)

	return nil
}
