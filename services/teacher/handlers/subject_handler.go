package handlers

import (
	"github.com/francoispqt/gojay"
	"github.com/juleur/dash-exp/apis/teacher/models"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

// CreateSubject func
func CreateSubject(tm *models.TeacherManager) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		subject := &models.Subject{}
		err := gojay.UnmarshalJSONObject(ctx.PostBody(), subject)
		if err != nil {
			tm.Log.WithFields(logrus.Fields{"value": subject}).Error(err)
			ctx.SetStatusCode(500)
			return
		}
		err = tm.DB.CreateSubject(*subject)
		if err != nil {
			tm.Log.Error(err)
			ctx.SetStatusCode(500)
			return
		}
		ctx.SetStatusCode(200)
	})
}
