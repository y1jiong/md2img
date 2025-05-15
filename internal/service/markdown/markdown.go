package markdown

import (
	_ "embed"
	"github.com/88250/lute"
	"log"
	"os"
	"strings"
)

const (
	placeholder = "{{.}}"
)

var (
	mathReplacer = strings.NewReplacer(
		`\[`, `\\[`,
		`\]`, `\\]`,
		`\(`, `\\(`,
		`\)`, `\\)`,
	)
	//go:embed template.html
	template string
)

func init() {
	contentBytes, err := os.ReadFile("template.html")
	if err != nil {
		return
	}
	content := string(contentBytes)
	if strings.Count(content, placeholder) != 1 {
		log.Println("template.html format error")
		return
	}
	template = content
	log.Println("custom template.html loaded")
}

// ToHTML 将Markdown转换为HTML
func ToHTML(md string) string {
	// 创建Lute引擎
	engine := lute.New()

	// 转换Markdown为HTML
	htmlContent := engine.MarkdownStr("", mathReplacer.Replace(md))

	// 处理Mermaid图表
	htmlContent = strings.ReplaceAll(htmlContent, `class="language-mermaid"`, `class="mermaid"`)

	// 构建完整的HTML文档
	return strings.Replace(template, placeholder, htmlContent, 1)
}
