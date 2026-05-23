package controller

import (
	"fmt"
	"io"
	"md2img/internal/service/browser"
	"net/http"
)

// URL 渲染端点
func URL(w http.ResponseWriter, r *http.Request) {
	// 读取请求体中的URL
	content, err := io.ReadAll(r.Body)
	if err != nil {
		sendError(w, http.StatusBadRequest, "无法读取请求体")
		return
	}
	defer r.Body.Close()

	targetURL := string(content)
	if targetURL == "" {
		sendError(w, http.StatusBadRequest, "URL不能为空")
		return
	}

	// Query parameters
	query := r.URL.Query()
	width, mobile, wait := queryWidth(query), queryMobile(query), queryWait(query)
	if width == 0 {
		if mobile {
			width = mobileWidth
		} else {
			width = desktopWidth
		}
	}

	// 渲染URL为图片
	imageData, err := browser.URL(targetURL, width, mobile, wait)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("渲染失败: %s", err))
		return
	}

	// 返回图片
	w.Header().Set("Content-Type", "image/png")
	_, _ = w.Write(imageData)
}
