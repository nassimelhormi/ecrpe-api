package handlers

import (
	"github.com/francoispqt/gojay"
	"github.com/juleur/dash-exp/apis/teacher/models"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

// CreateRefresherCourse func
func CreateRefresherCourse(tm *models.TeacherManager) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		rc := &models.RefresherCourse{}
		err := gojay.UnmarshalJSONObject(ctx.PostBody(), rc)
		if err != nil {
			tm.Log.WithFields(logrus.Fields{"value": rc}).Error(err)
			ctx.SetStatusCode(500)
			return
		}
		err = tm.DB.CreateRefresherCourse(*rc)
		if err != nil {
			tm.Log.Error(err)
			ctx.SetStatusCode(500)
			return
		}
		ctx.SetStatusCode(200)
	})
}
