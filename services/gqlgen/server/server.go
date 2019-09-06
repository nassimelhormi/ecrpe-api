package main

import (
	"net/http"

	"github.com/nassimelhormi/ecrpe-api/services/gqlgen/interceptors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"

	"github.com/go-chi/chi"

	"github.com/sirupsen/logrus"

	"github.com/99designs/gqlgen/handler"
	"github.com/nassimelhormi/ecrpe-api/services/gqlgen"
)

const (
	defaultPort = "8080"
	// env variable later
	secretKey = "secretkey"
)

var (
	db *sqlx.DB
)

func init() {
	db, err := sqlx.Open("mysql", "chermak:pwd@tcp(127.0.0.1:7359)/ecrpe")
	if err != nil {
		logrus.Fatalln(err)
	}
	if err = db.Ping(); err != nil {
		logrus.Fatalln(err)
	}
	defer db.Close()
}

func main() {
	router := chi.NewRouter()

	router.Use(interceptors.JWTCheck(secretKey))

	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:" + defaultPort},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	// [SECURITY] https://gqlgen.com/reference/complexity/
	// [APQ] https://gqlgen.com/reference/apq/
	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(gqlgen.NewExecutableSchema(gqlgen.Config{Resolvers: &gqlgen.Resolver{DB: db}})))

	logrus.Printf("connect to http://localhost:%s/ for GraphQL playground", defaultPort)
	if err := http.ListenAndServe(":"+defaultPort, nil); err != nil {
		logrus.Fatalln(err)
	}
}
