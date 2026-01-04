package markdown

import (
	"bytes"
	_ "embed"
	"log"
	"md2img/util"
	"os"
	"strings"
)

const (
	templateName = "template.html"
	mathjaxName  = "mathjax.html"
	mermaidName  = "mermaid.html"

	htmlPlaceholder    = "{{.}}"
	mathjaxPlaceholder = "{{mathjax}}"
	mermaidPlaceholder = "{{mermaid}}"

	blockMathL, blockMathR    = `\[`, `\]`
	inlineMathL, inlineMathR  = `\(`, `\)`
	blockMathMd, inlineMathMd = `$$`, `$`
)

var (
	//go:embed template.html
	templateHTML string
	//go:embed mathjax.html
	mathjaxHTML string
	//go:embed mermaid.html
	mermaidHTML string

	clearReplacer = strings.NewReplacer(
		mathjaxPlaceholder, "",
		mermaidPlaceholder, "",
	)
)

func init() {
	initTemplate()
	initTmpl(mathjaxName, &mathjaxHTML)
	initTmpl(mermaidName, &mermaidHTML)
}

// ToHTML 将Markdown转换为HTML
func ToHTML(md string, pure bool) string {
	// 初始化模板
	tmpl := templateHTML

	// 前置处理MathJax LaTeX
	if strings.Contains(md, blockMathL) || strings.Contains(md, inlineMathL) ||
		strings.Contains(md, blockMathMd) || strings.Contains(md, inlineMathMd) {
		tmpl = strings.Replace(tmpl, mathjaxPlaceholder, mathjaxHTML, 1)
	}

	// 转换Markdown为HTML
	var htmlContent string
	{
		var buf bytes.Buffer
		_ = engine.Convert(util.StringToBytes(md), &buf)
		htmlContent = buf.String()
	}

	// 后置处理Mermaid图表
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
	var content string
	initTmpl(templateName, &content)
	if content == "" {
		return
	}
	if strings.Count(content, htmlPlaceholder) != 1 {
		log.Println("custom", templateName, "invalid: missing or multiple", htmlPlaceholder)
		return
	}
	templateHTML = content
	log.Println("custom", templateName, "loaded")
}

func initTmpl(name string, target *string) {
	contentBytes, err := os.ReadFile(name)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Println("read", name, "failed:", err)
		return
	}
	*target = string(contentBytes)
	if name != templateName {
		log.Println("custom", name, "loaded")
	}
}
