package models

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

// TeacherManager struct
type TeacherManager struct {
	DB             *DB
	PublishVideoCh chan Video
	Log            *logrus.Logger
}

// NewTeacherManager struct
func NewTeacherManager() *TeacherManager {
	tm := &TeacherManager{
		DB:             NewDB(),
		PublishVideoCh: make(chan Video, 10),
		Log: &logrus.Logger{
			Formatter: &logrus.JSONFormatter{
				TimestampFormat: "02-01-2006 15:04:05.999999999",
				CallerPrettyfier: func(f *runtime.Frame) (string, string) {
					pwd, err := os.Getwd()
					if err != nil {
						logrus.Fatalln(err)
					}
					filename := strings.Replace(f.File, pwd, "", -1)
					return "", fmt.Sprintf("%s:%d", filename, f.Line)
				},
			},
			Level: logrus.ErrorLevel,
			Out: io.Writer(&lumberjack.Logger{
				Filename:   "./../../tools/logs/server.log",
				MaxSize:    5,
				MaxAge:     0,
				MaxBackups: 0,
			}),
			ReportCaller: true,
		},
	}
	return tm
}
