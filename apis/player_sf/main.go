package main

import (
	"github.com/AdhityaRamadhanus/fasthttpcors"
	"github.com/buaazp/fasthttprouter"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func main() {
	logrus.Infoln("Player Serve File Started")
	withCors := fasthttpcors.NewCorsHandler(fasthttpcors.Options{
		AllowedOrigins:   []string{},
		AllowedHeaders:   []string{},
		AllowedMethods:   []string{"GET"},
		AllowCredentials: false,
		Debug:            true,
	})

	router := fasthttprouter.New()
	router.ServeFiles("/api/*filepath", "/home/julien/player")
	if err := fasthttp.ListenAndServe(":4039", withCors.CorsMiddleware(router.Handler)); err != nil {
		logrus.Fatalln(err)
	}
}
