package middlewares

import (
	"context"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/dataloaders"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/repository"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/utils"
	"net/http"
)

func DataloaderMiddleware(repo repository.Repository, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loaders := dataloaders.NewDataLoaders(repo)
		ctx := context.WithValue(r.Context(), utils.DataLoadersCtxKey, loaders)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
