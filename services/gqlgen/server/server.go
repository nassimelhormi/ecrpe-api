package main

import (
	"net/http"
	"time"

	"github.com/nassimelhormi/ecrpe-api/services/gqlgen/directives"
	"github.com/nassimelhormi/ecrpe-api/services/gqlgen/redis"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"

	"github.com/go-chi/chi"

	"github.com/99designs/gqlgen/handler"
	"github.com/nassimelhormi/ecrpe-api/models"
	"github.com/nassimelhormi/ecrpe-api/services/gqlgen"
	"github.com/nassimelhormi/ecrpe-api/services/gqlgen/interceptors"
	"github.com/sirupsen/logrus"
)

const (
	defaultPort = "8080"
	secretKey   = "secretkey"
)

var (
	db              *sqlx.DB
	apqCache        *redis.Cache
	ipAddressCache  *redis.Cache
	videoEncodingCh chan models.Video
)

func init() {
	// secretKey = os.Getenv("SECRET_KET")
	// database
	videoEncodingCh = make(chan models.Video, 3)
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
					// 	DB:              db,
					// 	SecreyKey:       secretKey,
					// 	IPAddressCache:  ipAddressCache,
					// 	VideoEncodingCh: videoEncodingCh,
				},
				Directives: gqlgen.DirectiveRoot{
					HasRole:              directives.HasRole,
					RefresherCourseOwner: directives.RefresherCourseOwner,
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
