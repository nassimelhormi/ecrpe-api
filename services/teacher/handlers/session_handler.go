package handlers

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/juleur/dash-exp/apis/teacher/models"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

var (
	homePath = os.Getenv("HOME")
)

// dashManifestBuilder because running exec.Command won't work
func dashManifestBuilder(video models.Video) error {
	dashPath := video.FileName() + ".sh"
	f, err := os.OpenFile(dashPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	content := fmt.Sprintf(`#!/bin/bash
ffmpeg -f webm_dash_manifest -i %s_360p.webm -f webm_dash_manifest -i %s_480p.webm -f webm_dash_manifest -i %s_720p.webm -f webm_dash_manifest -i %s_audio.webm -c copy -map 0 -map 1 -map 2 -map 3 -f webm_dash_manifest -adaptation_sets "id=0,streams=0,1,2 id=1,streams=3" %s_manifest.mpd`, video.FileName(), video.FileName(), video.FileName(), video.FileName(), video.FileName())
	_, err = f.WriteString(content)
	return err
}

func mkdirPath(s models.Session) (string, error) {
	pathVideo := homePath + "/player/" + s.SubjectFolderName() + "/rcv/session-" + s.RefresherCourseYear
	if err := os.MkdirAll(pathVideo, 0777); err != nil {
		return "", err
	}
	return pathVideo, nil
}

func ffmpegProcessing(video models.Video, tm *models.TeacherManager) {
	argsAudio := []string{
		"-i", video.FileName(), "-c:a", "libopus", "-b:a", "64k", "-vn",
		"-f", "webm", "-dash", "1", video.FileName() + "_audio.webm",
	}
	argsVideo := []string{
		"-i", video.FileName(), "-c:v", "libvpx-vp9", "-keyint_min", "3", "-g", "30",
		"-speed", "2", "-f", "webm", "-dash", "1",
		"-an", "-vf", "scale=640x360", "-b:v", "276k", "-minrate", "138k", "-maxrate", "400k",
		"-quality", "good", "-crf", "45", "-dash", "1", video.FileName() + "_360p.webm",
		"-an", "-vf", "scale=640x480", "-b:v", "512k", "-minrate", "256k", "-maxrate", "742k",
		"-quality", "good", "-crf", "35", "-dash", "1", video.FileName() + "_480p.webm",
		"-an", "-vf", "scale=1280x720", "-b:v", "1024k", "-minrate", "512k", "-maxrate", "1485k",
		"-quality", "best", "-crf", "25", "-dash", "1", video.FileName() + "_720p.webm",
	}
	argsFfprobe := []string{
		"-v", "quiet", "-print_format", "compact=print_section=0:nokey=1:escape=csv", "-show_entries",
		"format=duration", video.FileName() + "_360p.webm",
	}
	cmd := exec.Command("ffmpeg", argsAudio...)
	if err := cmd.Run(); err != nil {
		tm.Log.Errorf("Audio encoding: %s", err)
		return
	}
	cmd = exec.Command("ffmpeg", argsVideo...)
	if err := cmd.Run(); err != nil {
		tm.Log.Errorf("Video encoding: %s", err)
		return
	}
	cmd = exec.Command("/bin/sh", video.FileName()+".sh")
	if err := cmd.Run(); err != nil {
		tm.Log.Errorf("Dash manifest: %s", err)
		return
	}
	duration, err := exec.Command("ffprobe", argsFfprobe...).Output()
	if err != nil {
		tm.Log.Errorf("Size video: %s", err)
		return
	}

	if err := os.Remove(video.FileName()); err != nil {
		tm.Log.Warnf("Remove %s file failed: %s", video.FileName(), err)
		return
	}
	if err := os.Remove(video.FileName() + ".sh"); err != nil {
		tm.Log.Warnf("Remove %s.sh failed: %s", video.FileName(), err)
		return
	}

	video.AddDuration(duration)
	video.Path = video.FileName() + "_manifest.mpd"
	video.IsEncoded = true

	tm.PublishVideoCh <- video
}

// CreateSession handler
func CreateSession(tm *models.TeacherManager) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		form, err := ctx.MultipartForm()
		if err != nil {
			tm.Log.WithFields(logrus.Fields{"value": form}).Error(err)
			ctx.SetStatusCode(500)
			return
		}
		session := &models.Session{}
		err = models.MultipartFormToSession(session, form)
		if err != nil {
			tm.Log.Error(err)
			ctx.SetStatusCode(500)
			return
		}
		video, err := tm.DB.CreateSessionVideo(*session)
		if err != nil {
			tm.Log.WithFields(logrus.Fields{"value": video}).Error(err)
			ctx.SetStatusCode(500)
			return
		}
		pathVideo, err := mkdirPath(*session)
		if err != nil {
			tm.Log.WithFields(logrus.Fields{"value": pathVideo}).Error(err)
			ctx.SetStatusCode(500)
			return
		}
		video.Path = pathVideo
		err = dashManifestBuilder(video)
		if err != nil {
			tm.Log.Error(err)
			ctx.SetStatusCode(500)
			return
		}
		videoFormField, err := ctx.FormFile("video")
		if err != nil {
			tm.Log.Error(err)
			ctx.SetStatusCode(500)
			return
		}
		err = fasthttp.SaveMultipartFile(videoFormField, video.FileName())
		if err != nil {
			tm.Log.Error(err)
			ctx.SetStatusCode(500)
			return
		}
		go ffmpegProcessing(video, tm)
		ctx.SetStatusCode(200)
	})
}
