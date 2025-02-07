package main

import (
	"OzonTestTask/internal/graphql"
	"OzonTestTask/internal/repository"
	"OzonTestTask/internal/service"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"log"
	"net/http"
)

const defaultPort = "8080"

func main() {
	repo := repository.NewInMemoryRepository()
	postService := service.NewPostService(repo)

	handler.New(graphql.NewExecutableSchema(graphql.Config{
		Resolvers: &graphql.Resolver{PostService: postService},
	}))
	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{
		Resolvers: &graphql.Resolver{
			PostService: postService,
		},
	}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", defaultPort)
	log.Fatal(http.ListenAndServe(":"+defaultPort, nil))
}
