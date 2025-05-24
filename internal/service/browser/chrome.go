package browser

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	fileUrlPrefix = "file://"

	contextTimeout = 30 * time.Second
)

var (
	allocatorCtx context.Context
	browserCtx   context.Context

	initBrowser = sync.OnceFunc(func() {
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
)

func HTML(html []byte, width int64, mobile bool) ([]byte, error) {
	// 创建临时HTML文件
	tempFile, err := os.CreateTemp("", "render-*.html")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write(html); err != nil {
		return nil, err
	}
	tempFile.Close()

	// 渲染为图片
	return URL(fileUrlPrefix+tempFile.Name(), width, mobile)
}

func URL(url string, width int64, mobile bool) ([]byte, error) {
	initBrowser()

	tabCtx, tabCancel := chromedp.NewContext(browserCtx)
	defer tabCancel()

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(tabCtx, contextTimeout)
	defer cancel()

	// 运行Chrome任务
	if err := chromedp.Run(ctx, chromedp.Tasks{
		emulation.SetDeviceMetricsOverride(width, 0, 2, mobile),
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
	}); err != nil {
		return nil, fmt.Errorf("浏览器渲染失败: %w", err)
	}

	{
		// 等待页面加载完成
		var readyState string
		startTime := time.Now()
		for time.Since(startTime) < 10*time.Second {
			if err := chromedp.Run(ctx, chromedp.Evaluate(`document.readyState`, &readyState)); err != nil {
				return nil, fmt.Errorf("获取页面状态失败: %w", err)
			}
			if readyState == "complete" {
				break
			}
			// 等待页面加载完成
			time.Sleep(200 * time.Millisecond)
		}
	}

	{
		// 检查并等待MathJax渲染完成
		var mathjaxRendered bool
		startTime := time.Now()

		for time.Since(startTime) < 10*time.Second {
			if err := chromedp.Run(ctx, chromedp.Evaluate(`
(async () => {
	try {
		await window.MathJax?.typesetPromise?.();
		return true;
	} catch (e) {
		return false;
	}
})()`,
				&mathjaxRendered,
				func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
					p = p.WithTimeout(10 * 1000)
					return p.WithAwaitPromise(true)
				},
			)); err != nil {
				return nil, fmt.Errorf("检查/等待MathJax失败: %w", err)
			}
			if mathjaxRendered {
				time.Sleep(100 * time.Millisecond)
				break
			}
			// 等待MathJax渲染完成
			time.Sleep(200 * time.Millisecond)
		}
	}

	if !strings.HasPrefix(url, fileUrlPrefix) {
		// 等待页面稳定
		time.Sleep(500 * time.Millisecond)
	}

	// 用于存储屏幕截图的字节数组
	var buf []byte

	// 执行截图任务
	if err := chromedp.Run(ctx, chromedp.FullScreenshot(&buf, 100)); err != nil {
		return nil, fmt.Errorf("截图失败: %w", err)
	}

	return buf, nil
}
