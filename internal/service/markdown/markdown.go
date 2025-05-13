package markdown

import (
	_ "embed"
	"fmt"
	"github.com/88250/lute"
	"log"
	"os"
	"strings"
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
	if strings.Count(content, "%s") != 1 {
		log.Println("template.html format error")
		return
	}
	template = content
	log.Println("custom template.html loaded")
}

// ToHTML 将Markdown转换为HTML
func ToHTML(md string) string {
	md = mathReplacer.Replace(md)

	// 创建Lute引擎
	engine := lute.New()

	// 转换Markdown为HTML
	htmlContent := engine.MarkdownStr("", md)

	// 构建完整的HTML文档
	return fmt.Sprintf(template, htmlContent)
}
