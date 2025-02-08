package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/graphql"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/middleware"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/repository"
	"github.com/ArtemSarafannikov/OzonTestTask/internal/service"
	"log"
	"net/http"
)

const defaultPort = "8080"

func main() {
	repo := repository.NewInMemoryRepository()
	postService := service.NewPostService(repo)
	userService := service.NewUserService(repo)
	commentService := service.NewCommentService(repo)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.NewRootResolver(postService, userService, commentService)))

	http.Handle("/query", middleware.AuthMiddleware(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", defaultPort)
	log.Fatal(http.ListenAndServe(":"+defaultPort, nil))
}
