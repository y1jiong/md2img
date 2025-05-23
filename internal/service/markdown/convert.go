package markdown

import (
	"github.com/gohugoio/hugo-goldmark-extensions/extras"
	"github.com/gohugoio/hugo-goldmark-extensions/passthrough"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
	"reflect"
)

var (
	engine = goldmark.New(
		goldmark.WithExtensions(
			extension.CJK,
			extension.Table,    // Part of GFM
			extension.Linkify,  // Part of GFM
			extension.TaskList, // Part of GFM
			extension.Footnote,
			highlighting.Highlighting,
			passthrough.New(
				passthrough.Config{
					BlockDelimiters: []passthrough.Delimiters{
						{
							Open:  blockMathMd,
							Close: blockMathMd,
						},
						{
							Open:  blockMathL,
							Close: blockMathR,
						},
					},
					InlineDelimiters: []passthrough.Delimiters{
						{
							Open:  inlineMathMd,
							Close: inlineMathMd,
						},
						{
							Open:  inlineMathL,
							Close: inlineMathR,
						},
					},
				},
			),
			extras.New(
				extras.Config{
					Delete:      extras.DeleteConfig{Enable: true},
					Insert:      extras.InsertConfig{Enable: true},
					Mark:        extras.MarkConfig{Enable: true},
					Subscript:   extras.SubscriptConfig{Enable: true},
					Superscript: extras.SuperscriptConfig{Enable: true},
				},
			),
		),
		goldmark.WithParserOptions(
			withoutSetextHeadingParser{},
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
)

type withoutSetextHeadingParser struct{}

func (w withoutSetextHeadingParser) SetParserOption(c *parser.Config) {
	setextType := reflect.TypeOf(parser.NewSetextHeadingParser())

	for _, v := range c.BlockParsers {
		if reflect.TypeOf(v.Value) != setextType {
			continue
		}
		c.BlockParsers = append(util.PrioritizedSlice{}, c.BlockParsers.Remove(v.Value)...)
		break
	}
}
