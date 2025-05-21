package markdown

import (
	"bytes"
	_ "embed"
	"github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"log"
	"md2img/util"
	"os"
	"strings"
)

const (
	htmlPlaceholder    = "{{.}}"
	mathjaxPlaceholder = "{{mathjax}}"
	mermaidPlaceholder = "{{mermaid}}"

	blockMath  = `\[`
	inlineMath = `\(`
)

var (
	//go:embed template.html
	templateHTML string
	//go:embed mathjax.html
	mathjaxHTML string
	//go:embed mermaid.html
	mermaidHTML string

	mathReplacer = strings.NewReplacer(
		`\[`, `$$`,
		`\]`, `$$`,
		`\(`, `$`,
		`\)`, `$`,
	)
	clearReplacer = strings.NewReplacer(
		mathjaxPlaceholder, "",
		mermaidPlaceholder, "",
	)

	engine = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.CJK,
			extension.Footnote,
			highlighting.Highlighting,
			mathjax.MathJax,
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
)

func init() {
	initTemplate()
	initMathjax()
	initMermaid()
}

// ToHTML 将Markdown转换为HTML
func ToHTML(md string, pure bool) string {
	// 初始化模板
	tmpl := templateHTML

	// 预处理MathJax LaTeX
	mathjaxLoaded := false
	if strings.Contains(md, blockMath) || strings.Contains(md, inlineMath) {
		tmpl = strings.Replace(tmpl, mathjaxPlaceholder, mathjaxHTML, 1)
		md = mathReplacer.Replace(md)
		mathjaxLoaded = true
	}

	// 转换Markdown为HTML
	var htmlContent string
	{
		var buf bytes.Buffer
		_ = engine.Convert(util.StringToBytes(md), &buf)
		htmlContent = buf.String()
	}

	// 处理MathJax LaTeX
	if !mathjaxLoaded && (strings.Contains(htmlContent, blockMath) || strings.Contains(htmlContent, inlineMath)) {
		tmpl = strings.Replace(tmpl, mathjaxPlaceholder, mathjaxHTML, 1)
	}

	// 处理Mermaid图表
	if strings.Contains(htmlContent, `class="language-mermaid"`) {
		tmpl = strings.Replace(tmpl, mermaidPlaceholder, mermaidHTML, 1)
		htmlContent = strings.ReplaceAll(htmlContent, `class="language-mermaid"`, `class="mermaid"`)
	}

	if pure {
		return htmlContent
	}

	// 清理占位符
	tmpl = clearReplacer.Replace(tmpl)

	// 构建完整的HTML文档
	return strings.Replace(tmpl, htmlPlaceholder, htmlContent, 1)
}

func initTemplate() {
	contentBytes, err := os.ReadFile("template.html")
	if err != nil {
		return
	}
	content := string(contentBytes)
	if strings.Count(content, htmlPlaceholder) != 1 {
		log.Println("template.html format error")
		return
	}
	templateHTML = content
	log.Println("custom template.html loaded")
}

func initMathjax() {
	contentBytes, err := os.ReadFile("mathjax.html")
	if err != nil {
		return
	}
	mathjaxHTML = string(contentBytes)
	log.Println("custom mathjax.html loaded")
}

func initMermaid() {
	contentBytes, err := os.ReadFile("mermaid.html")
	if err != nil {
		return
	}
	mermaidHTML = string(contentBytes)
	log.Println("custom mermaid.html loaded")
}
