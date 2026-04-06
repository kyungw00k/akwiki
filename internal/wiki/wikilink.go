package wiki

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// Wikilink AST 노드
type Wikilink struct {
	ast.BaseInline
	Target  string
	Display string
}

var KindWikilink = ast.NewNodeKind("Wikilink")

func (n *Wikilink) Kind() ast.NodeKind { return KindWikilink }

func (n *Wikilink) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{
		"Target":  n.Target,
		"Display": n.Display,
	}, nil)
}

// wikilinkParser는 [[...]] 패턴을 파싱합니다.
type wikilinkParser struct{}

var defaultWikilinkParser = &wikilinkParser{}

func (p *wikilinkParser) Trigger() []byte {
	return []byte{'['}
}

func (p *wikilinkParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()

	// [[ 로 시작하는지 확인
	if len(line) < 4 || line[0] != '[' || line[1] != '[' {
		return nil
	}

	// ]] 찾기
	rest := string(line[2:])
	end := strings.Index(rest, "]]")
	if end < 0 {
		return nil
	}

	inner := rest[:end]
	target := inner
	display := inner

	// | 로 target / display 분리
	if idx := strings.Index(inner, "|"); idx >= 0 {
		target = inner[:idx]
		display = inner[idx+1:]
	}

	// 리더를 [[ + inner + ]] 만큼 전진
	block.Advance(2 + end + 2)

	node := &Wikilink{
		Target:  target,
		Display: display,
	}
	return node
}

// wikilinkRenderer는 Wikilink 노드를 HTML로 렌더링합니다.
type wikilinkRenderer struct {
	pageRoute string
}

func (r *wikilinkRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindWikilink, r.renderWikilink)
}

func (r *wikilinkRenderer) renderWikilink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*Wikilink)
	href := r.pageRoute + "/" + url.PathEscape(n.Target)
	fmt.Fprintf(w, `<a href="%s" class="wikilink">%s</a>`, href, n.Display)
	return ast.WalkContinue, nil
}

// wikilinkExtension은 goldmark 확장입니다.
type wikilinkExtension struct {
	pageRoute string
}

func NewWikilinkExtension(pageRoute string) goldmark.Extender {
	return &wikilinkExtension{pageRoute: pageRoute}
}

func (e *wikilinkExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(defaultWikilinkParser, 150),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&wikilinkRenderer{pageRoute: e.pageRoute}, 200),
		),
	)
}

// LinkStatus represents the resolution status of a wikilink target.
type LinkStatus int

const (
	LinkExists  LinkStatus = iota
	LinkPrivate
	LinkMissing
)

// LinkResolver determines the status of a wikilink target.
type LinkResolver func(target string) LinkStatus

// wikilinkResolveRenderer renders Wikilink nodes with link resolution.
type wikilinkResolveRenderer struct {
	pageRoute string
	resolver  LinkResolver
}

func (r *wikilinkResolveRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindWikilink, r.renderWikilink)
}

func (r *wikilinkResolveRenderer) renderWikilink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*Wikilink)

	status := LinkExists
	if r.resolver != nil {
		status = r.resolver(n.Target)
	}

	switch status {
	case LinkPrivate:
		fmt.Fprintf(w, `<a href="#private-link" data-link="private" class="wikilink">%s</a>`, n.Display)
	case LinkMissing:
		fmt.Fprintf(w, `<a class="wikilink-missing">%s</a>`, n.Display)
	default:
		href := r.pageRoute + "/" + url.PathEscape(n.Target)
		fmt.Fprintf(w, `<a href="%s" class="wikilink">%s</a>`, href, n.Display)
	}
	return ast.WalkContinue, nil
}

// wikilinkResolveExtension is a goldmark extension with link resolution.
type wikilinkResolveExtension struct {
	pageRoute string
	resolver  LinkResolver
}

// NewWikilinkExtensionWithResolver creates a wikilink extension that uses a resolver
// to determine the status of each link target.
func NewWikilinkExtensionWithResolver(pageRoute string, resolver LinkResolver) goldmark.Extender {
	return &wikilinkResolveExtension{pageRoute: pageRoute, resolver: resolver}
}

func (e *wikilinkResolveExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(defaultWikilinkParser, 150),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&wikilinkResolveRenderer{
				pageRoute: e.pageRoute,
				resolver:  e.resolver,
			}, 200),
		),
	)
}
