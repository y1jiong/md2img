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

	url := string(content)
	if url == "" {
		sendError(w, http.StatusBadRequest, "URL不能为空")
		return
	}

	// 渲染URL为图片
	imageData, err := browser.URL(url)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("渲染失败: %s", err))
		return
	}

	// 返回图片
	_, _ = w.Write(imageData)
}
