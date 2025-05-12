package browser

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"os"
	"sync"
	"time"
)

const (
	contextTimeout = 30 * time.Second
)

var (
	allocatorCtx context.Context
	browserCtx   context.Context
	initOnce     sync.Once
)

func initBrowser() {
	initOnce.Do(func() {
		// 创建无头Chrome实例
		options := []chromedp.ExecAllocatorOption{
			chromedp.Flag("headless", "new"),          // 启用全新无头模式
			chromedp.Flag("disable-gpu", true),        // 避免 GPU 渲染问题
			chromedp.Flag("no-sandbox", true),         // 有些系统必须禁用 sandbox
			chromedp.Flag("enable-automation", false), // 禁用自动化提示
		}
		allocatorCtx, _ = chromedp.NewExecAllocator(
			context.Background(),
			append(chromedp.DefaultExecAllocatorOptions[:], options...)...,
		)
		browserCtx, _ = chromedp.NewContext(allocatorCtx)
	})
}

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
	initBrowser()

	tabCtx, tabCancel := chromedp.NewContext(browserCtx)
	defer tabCancel()

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(tabCtx, contextTimeout)
	defer cancel()

	// 运行Chrome任务
	if err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
	}); err != nil {
		return nil, fmt.Errorf("浏览器渲染失败: %w", err)
	}

	// 等待页面加载完成
	var readyState string
	for readyState != "complete" {
		time.Sleep(100 * time.Millisecond) // 避免过于频繁的检查
		if err := chromedp.Run(ctx, chromedp.Evaluate("document.readyState", &readyState)); err != nil {
			return nil, fmt.Errorf("获取页面状态失败: %w", err)
		}
	}

	// 用于存储屏幕截图的字节数组
	var buf []byte

	// 执行截图任务
	if err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Sleep(500 * time.Millisecond), // 等待页面稳定
		chromedp.FullScreenshot(&buf, 100),
	}); err != nil {
		return nil, fmt.Errorf("截图失败: %w", err)
	}

	return buf, nil
}
