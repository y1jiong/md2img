package browser

import (
	"context"
	"errors"
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
	browserCtx context.Context

	initBrowser = sync.OnceFunc(func() {
		// 创建无头Chrome实例
		options := []chromedp.ExecAllocatorOption{
			chromedp.Flag("headless", "new"),          // 启用全新无头模式
			chromedp.Flag("disable-gpu", true),        // 避免 GPU 渲染问题
			chromedp.Flag("no-sandbox", true),         // 有些系统必须禁用 sandbox
			chromedp.Flag("enable-automation", false), // 禁用自动化提示
		}
		allocatorCtx, _ := chromedp.NewExecAllocator(
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
		chromedp.Poll(`document.readyState === 'complete'`, nil),
	}); err != nil {
		return nil, fmt.Errorf("浏览器渲染失败: %w", err)
	}

	if !strings.HasPrefix(url, fileUrlPrefix) {
		if err := chromedp.Run(ctx, waitForDOMStable(150, time.Second)); err != nil {
			return nil, fmt.Errorf("DOM稳定性检查失败: %w", err)
		}
	}

	// 用于存储屏幕截图的字节数组
	var buf []byte

	// 执行截图任务
	if err := chromedp.Run(ctx, chromedp.FullScreenshot(&buf, 100)); err != nil {
		return nil, fmt.Errorf("截图失败: %w", err)
	}

	return buf, nil
}

func waitForDOMStable(settleDelayMS int, timeout time.Duration) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		const script = `
new Promise(resolve => {
	let timeout;
	const observer = new MutationObserver(() => {
		clearTimeout(timeout);
		timeout = setTimeout(() => {
			observer.disconnect();
			resolve();
		}, %d);
	});
	observer.observe(document.body, { childList: true, subtree: true });

	// 初始 fallback，确保 resolve 一定发生
	timeout = setTimeout(() => {
		observer.disconnect();
		resolve();
	}, %d);
})`
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		if err := chromedp.Evaluate(
			fmt.Sprintf(script, settleDelayMS, settleDelayMS),
			nil,
			func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
				return p.WithAwaitPromise(true)
			},
		).Do(ctx); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		return nil
	}
}
