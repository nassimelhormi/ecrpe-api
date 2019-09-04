package main

import (
	"net/http"

	"github.com/rs/cors"

	"github.com/go-chi/chi"

	"github.com/sirupsen/logrus"

	"github.com/99designs/gqlgen/handler"
	"github.com/nassimelhormi/ecrpe-api/services/gqlgen"
)

const defaultPort = "8080"

func main() {
	router := chi.NewRouter()

	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:" + defaultPort},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	// [SECURITY] https://gqlgen.com/reference/complexity/
	// [APQ] https://gqlgen.com/reference/apq/
	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(gqlgen.NewExecutableSchema(gqlgen.Config{Resolvers: &gqlgen.Resolver{}})))

	logrus.Printf("connect to http://localhost:%s/ for GraphQL playground", defaultPort)
	if err := http.ListenAndServe(":"+defaultPort, nil); err != nil {
		logrus.Fatalln(err)
	}
}
