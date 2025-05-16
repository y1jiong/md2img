package controller

import (
	"fmt"
	"github.com/88250/lute/util"
	"io"
	"md2img/internal/service/browser"
	"md2img/internal/service/markdown"
	"net/http"
)

// Markdown 渲染端点
func Markdown(w http.ResponseWriter, r *http.Request) {
	// 读取请求体中的Markdown内容
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
		width = mobileWidth
	}

	// 转换Markdown为HTML
	html := markdown.ToHTML(util.BytesToStr(content))

	// 渲染为图片
	imageData, err := browser.HTML(html, width, mobile)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("渲染失败: %s", err))
		return
	}

	// 返回图片
	_, _ = w.Write(imageData)
}
