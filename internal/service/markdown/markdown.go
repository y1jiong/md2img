package markdown

import (
	"github.com/88250/lute"
	"strings"
)

var (
	mathReplacer = strings.NewReplacer(
		`\[`, `\\[`,
		`\]`, `\\]`,
		`\(`, `\\(`,
		`\)`, `\\)`,
	)
	front = `<!DOCTYPE html>
<html lang="zh">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <!-- GitHub Markdown Light Style -->
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/github-markdown-css@5.8.1/github-markdown-light.min.css">
  <style>
    body {
      background-color: white;
      display: flex;
      justify-content: center;
      padding: 2rem;
    }
    .markdown-body {
      max-width: 800px;
    }
  </style>
  <!-- MathJax for rendering LaTeX -->
  <script>
    window.MathJax = {
      tex: { inlineMath: [['\\(', '\\)']], displayMath: [['\\[', '\\]']] },
      svg: { fontCache: 'global' }
    };
  </script>
  <script async src="https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-mml-chtml.js"></script>
</head>
<body>
  <article class="markdown-body">
`
	end = `
  </article>
</body>
</html>`
)

// ToHTML 将Markdown转换为HTML
func ToHTML(md string) string {
	md = mathReplacer.Replace(md)

	// 创建Lute引擎
	engine := lute.New()

	// 转换Markdown为HTML
	htmlContent := engine.MarkdownStr("", md)

	// 构建完整的HTML文档
	return front + htmlContent + end
}
