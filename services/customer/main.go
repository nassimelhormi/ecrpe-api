package main

import (
	"github.com/savsgio/atreugo"
	"github.com/sirupsen/logrus"
)

func main() {
	config := &atreugo.Config{
		Host:             "0.0.0.0",
		Port:             5858,
		TLSEnable:        false,
		LogLevel:         "debug",
		Compress:         true,
		GracefulShutdown: true,
		Fasthttp: &atreugo.FasthttpConfig{
			Name:               "Customer Service",
			MaxRequestBodySize: 200 << 20, // 200MB
		},
	}
	server := atreugo.New(config)

	/// add CORS to test with graphQL
	/// add paths there

	server.Path("POST", "/query", viewFn atreugo.View)
	
	if err := server.ListenAndServe(); err != nil {
		logrus.Fatalln(err)
	}
}
