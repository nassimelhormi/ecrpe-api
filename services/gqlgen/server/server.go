package main

import (
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"

	"github.com/go-chi/chi"

	"github.com/sirupsen/logrus"

	"github.com/99designs/gqlgen/handler"
	"github.com/nassimelhormi/ecrpe-api/services/gqlgen"
	"github.com/nassimelhormi/ecrpe-api/services/gqlgen/apq"
	"github.com/nassimelhormi/ecrpe-api/services/gqlgen/interceptors"
)

const (
	defaultPort = "8080"
	secretKey   = "secretkey"
)

var (
	db       *sqlx.DB
	apqCache *apq.Cache
)

func init() {
	// secretKey = os.Getenv("SECRET_KET")
	db, err := sqlx.Open("mysql", "chermak:pwd@tcp(127.0.0.1:7359)/ecrpe")
	if err != nil {
		logrus.Fatalln(err)
	}
	if err = db.Ping(); err != nil {
		logrus.Fatalln(err)
	}
	defer db.Close()

	apqCache, err = apq.NewCache("localhost:4567", "", 24*time.Hour)
	if err != nil {
		logrus.Fatalf("cannot create APQ redis cache: %v", err)
	}
}

func main() {
	router := chi.NewRouter()

	router.Use(interceptors.JWTCheck(secretKey))
	router.Use(interceptors.GetIPAddress())

	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:" + defaultPort},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(
		gqlgen.NewExecutableSchema(
			gqlgen.Config{
				Resolvers: &gqlgen.Resolver{
					DB:        db,
					SecreyKey: secretKey,
				},
			},
		),
		handler.ComplexityLimit(3),
		handler.EnablePersistedQueryCache(apqCache),
	))

	logrus.Printf("connect to http://localhost:%s/ for GraphQL playground\n", defaultPort)
	if err := http.ListenAndServe(":"+defaultPort, nil); err != nil {
		logrus.Fatalln(err)
	}
}
