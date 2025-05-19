package controller

import (
	"fmt"
	"io"
	"md2img/internal/service/browser"
	"net/http"
)

// HTML 渲染端点
func HTML(w http.ResponseWriter, r *http.Request) {
	// 读取请求体中的HTML内容
	content, err := io.ReadAll(r.Body)
	if err != nil {
		sendError(w, http.StatusBadRequest, "无法读取请求体")
		return
	}
	defer r.Body.Close()

	// Query parameters
	query := r.URL.Query()
	width, mobile := queryWidth(query), queryMobile(query)
	if width == 0 {
		if mobile {
			width = mobileWidth
		} else {
			width = desktopWidth
		}
	}

	// 渲染HTML为图片
	imageData, err := browser.HTML(content, width, mobile)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("渲染失败: %s", err))
		return
	}

	// 返回图片
	_, _ = w.Write(imageData)
}
