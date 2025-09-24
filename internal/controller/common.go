package controller

import (
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	mobileWidth  = 440
	desktopWidth = 1280
)

func queryWidth(query url.Values) (width int64) {
	widthStr := query.Get("width")
	if widthStr == "" {
		return
	}
	width, _ = strconv.ParseInt(widthStr, 10, 64)
	return
}

func queryMobile(query url.Values) (mobile bool) {
	switch query.Get("mobile") {
	case "1", "true", "True", "TRUE":
		mobile = true
	}
	return
}

func queryWait(query url.Values) (wait time.Duration) {
	waitStr := query.Get("wait")
	if waitStr == "" {
		return
	}
	wait, _ = time.ParseDuration(waitStr)
	return
}

// 发送错误响应
func sendError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	_, _ = w.Write([]byte(message))
}
