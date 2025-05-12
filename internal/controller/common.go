package controller

import "net/http"

// 发送错误响应
func sendError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	_, _ = w.Write([]byte(message))
}
