package main

import (
	"net/http"
	"time"

	"github.com/plutov/paypal"

	"github.com/nassimelhormi/ecrpe-api/services/gqlgen/redis"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"

	"github.com/go-chi/chi"

	"github.com/sirupsen/logrus"

	"github.com/99designs/gqlgen/handler"
	"github.com/nassimelhormi/ecrpe-api/services/gqlgen"
	"github.com/nassimelhormi/ecrpe-api/services/gqlgen/interceptors"
)

const (
	defaultPort    = "8080"
	secretKey      = "secretkey"
	paypalClientID = "Af4joGyOARBhHtZz0LFQhZiXHdJV_Xiaih97aftQqrPCRowR9dO7WuoN9CB9ZynaJzbPn9iWlZTDHw1n"
	paypalSecretID = "EE5erT0UHU0-kn1O3X-gW3Xrj2ddIgMy8qNvrVPZ7eOSQHOZ0YIl1fOK8xzm8dJ7hA9ZHhn5lC_7UfBj"
)

var (
	db             *sqlx.DB
	apqCache       *redis.Cache
	ipAddressCache *redis.Cache
	paypalClient   *paypal.Client
)

func init() {
	// secretKey = os.Getenv("SECRET_KET")
	// paypalClientID = os.Getenv("PAYPAL_CLIENT_ID")
	// paypalSecretID = os.Getenv("SECRET_CLIENT_ID")
	// database
	db, err := sqlx.Open("mysql", "chermak:pwd@tcp(127.0.0.1:7359)/ecrpe")
	if err != nil {
		logrus.Fatalf("cannot connect to mysql: %v", err)
	}
	if err := db.Ping(); err != nil {
		logrus.Fatalln("database ping didn't work", err)
	}
	defer db.Close()
	// redis servers
	if apqCache, err = redis.NewCache("localhost:4567", "", 24*time.Hour); err != nil {
		logrus.Fatalf("cannot create APQ redis cache: %v", err)
	}
	if ipAddressCache, err = redis.NewCache("localhost:6789", "", 15*time.Minute); err != nil {
		logrus.Fatalf("cannot create IP Address redis cache: %v", err)
	}
	// paypal client
	if paypalClient, err = paypal.NewClient(paypalClientID, paypalSecretID, paypal.APIBaseSandBox); err != nil {
		logrus.Fatalf("cannot connect to paypal service: %v", err)
	}
}

func main() {
	router := chi.NewRouter()

	router.Use(interceptors.JWTCheck(secretKey))
	router.Use(interceptors.GetIPAddress())
	router.Use(interceptors.GetRefreshToken())

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
					DB:             db,
					SecreyKey:      secretKey,
					IPAddressCache: ipAddressCache,
					PaypalClient:   paypalClient,
				},
			},
		),
		handler.UploadMaxSize(300),
		handler.UploadMaxMemory(300),
		handler.ComplexityLimit(3),
		handler.EnablePersistedQueryCache(apqCache),
	))

	logrus.Printf("connect to http://localhost:%s/ for GraphQL playground\n", defaultPort)
	if err := http.ListenAndServe(":"+defaultPort, nil); err != nil {
		logrus.Fatalln(err)
	}
}
