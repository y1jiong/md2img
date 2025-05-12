package browser

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"os"
	"time"
)

const (
	contextTimeout = 30 * time.Second
)

func HTML(html string) ([]byte, error) {
	// 创建临时HTML文件
	tempFile, err := os.CreateTemp("", "render-*.html")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write([]byte(html)); err != nil {
		return nil, err
	}
	tempFile.Close()

	// 渲染为图片
	return URL("file://" + tempFile.Name())
}

func URL(url string) ([]byte, error) {
	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(
		context.Background(),
		contextTimeout,
	)
	defer cancel()

	// 创建无头Chrome实例
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),                      // 启用无头模式（可设为 false 调试用）
		chromedp.Flag("disable-gpu", true),                   // 避免 GPU 渲染问题
		chromedp.Flag("no-sandbox", true),                    // 有些系统必须禁用 sandbox
		chromedp.Flag("disable-dev-shm-usage", true),         // 避免 /dev/shm 空间不足的问题（Docker 环境尤为常见）
		chromedp.Flag("disable-extensions", true),            // 提升性能，避免加载扩展
		chromedp.Flag("hide-scrollbars", true),               // 隐藏滚动条
		chromedp.Flag("mute-audio", true),                    // 静音，提升性能
		chromedp.Flag("disable-background-networking", true), // 禁用后台网络请求
		chromedp.Flag("disable-default-apps", true),          // 禁用默认应用
		chromedp.Flag("disable-sync", true),                  // 禁用数据同步
	}

	actx, acancel := chromedp.NewExecAllocator(ctx, options...)
	defer acancel()

	// 创建Chrome上下文
	browserCtx, bcancel := chromedp.NewContext(actx)
	defer bcancel()

	// 用于存储屏幕截图的字节数组
	var buf []byte

	// 运行Chrome任务
	if err := chromedp.Run(browserCtx, screenshotTasks(url, &buf)); err != nil {
		return nil, fmt.Errorf("浏览器渲染失败: %w", err)
	}

	return buf, nil
}

// 截图任务
func screenshotTasks(url string, buf *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		//chromedp.Evaluate(`MathJax.typesetPromise().then(() => true)`, nil),
		chromedp.Sleep(2 * time.Second), // 给一点缓冲时间
		chromedp.FullScreenshot(buf, 100),
	}
}
