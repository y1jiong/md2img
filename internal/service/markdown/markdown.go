package markdown

import (
	_ "embed"
	"github.com/88250/lute"
	"log"
	"os"
	"strings"
)

const (
	htmlPlaceholder    = "{{.}}"
	mathjaxPlaceholder = "{{mathjax}}"
	mermaidPlaceholder = "{{mermaid}}"
)

var (
	//go:embed template.html
	template string
	//go:embed mathjax.html
	mathjax string
	//go:embed mermaid.html
	mermaid string

	mathReplacer = strings.NewReplacer(
		`\[`, `\\[`,
		`\]`, `\\]`,
		`\(`, `\\(`,
		`\)`, `\\)`,
	)
	clearReplacer = strings.NewReplacer(
		mathjaxPlaceholder, "",
		mermaidPlaceholder, "",
	)
)

func init() {
	initTemplate()
	initMathjax()
	initMermaid()
}

// ToHTML 将Markdown转换为HTML
func ToHTML(md string) string {
	// 初始化模板
	tmpl := template

	// 创建Lute引擎
	engine := lute.New()

	// 处理MathJax LaTeX
	if strings.Contains(md, `\[`) || strings.Contains(md, `\(`) {
		tmpl = strings.Replace(tmpl, mathjaxPlaceholder, mathjax, 1)
		md = mathReplacer.Replace(md)
	}

	// 转换Markdown为HTML
	htmlContent := engine.MarkdownStr("", md)

	// 处理Mermaid图表
	if strings.Contains(htmlContent, `class="language-mermaid"`) {
		tmpl = strings.Replace(tmpl, mermaidPlaceholder, mermaid, 1)
		htmlContent = strings.ReplaceAll(htmlContent, `class="language-mermaid"`, `class="mermaid"`)
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
	template = content
	log.Println("custom template.html loaded")
}

func initMathjax() {
	contentBytes, err := os.ReadFile("mathjax.html")
	if err != nil {
		return
	}
	mathjax = string(contentBytes)
	log.Println("custom mathjax.html loaded")
}

func initMermaid() {
	contentBytes, err := os.ReadFile("mermaid.html")
	if err != nil {
		return
	}
	mermaid = string(contentBytes)
	log.Println("custom mermaid.html loaded")
}
