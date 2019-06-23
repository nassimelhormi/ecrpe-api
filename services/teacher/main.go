package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/juleur/dash-exp/apis/teacher/handlers"
	"github.com/juleur/dash-exp/apis/teacher/models"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func main() {
	logrus.Infoln("Teacher Service Started")

	tm := models.NewTeacherManager()
	defer tm.DB.Close()

	go models.PublisherVideo(tm)

	router := fasthttprouter.New()
	router.POST("/create_subject", handlers.CreateSubject(tm))
	router.POST("/create_refresher_course", handlers.CreateRefresherCourse(tm))
	router.POST("/create_session", handlers.CreateSession(tm))
	s := &fasthttp.Server{
		Name:               "Teacher Service",
		Handler:            router.Handler,
		MaxRequestBodySize: 200 << 20, // 200MB
	}
	if err := s.ListenAndServe("127.0.0.1:5656"); err != nil {
		logrus.Fatalln(err)
	}
}
