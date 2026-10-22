// Package htmlfmt indents HTML so "view source" is readable. templ emits
// compact markup, and this re-flows it without changing how it renders:
// elements that directly contain text (and everything inside them) stay
// inline, whitespace-sensitive elements are left verbatim, and only
// structural nesting gets line breaks.
package htmlfmt

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

const indent = "  "

var voidElements = map[string]bool{
	"area": true, "base": true, "br": true, "col": true, "embed": true, "hr": true,
	"img": true, "input": true, "link": true, "meta": true, "source": true,
	"track": true, "wbr": true,
}

// verbatim elements keep their content exactly as written.
var verbatim = map[string]bool{"pre": true, "textarea": true, "script": true, "style": true}

type token struct {
	typ   html.TokenType
	raw   string
	name  string
	mixed bool // start tag whose content must stay inline
}

type frame struct {
	name   string
	mixed  bool
	inFlow bool // opened inside a mixed element: no indentation bookkeeping
}

// Format returns src re-indented. If src cannot be tokenized it is returned unchanged.
func Format(src []byte) []byte {
	toks, ok := tokenize(src)
	if !ok {
		return src
	}

	var (
		b        bytes.Buffer
		stack    []frame
		depth    int // indent level of non-flow elements
		flow     int // number of open mixed elements
		lastOpen = -1
	)
	breakLine := func() {
		if b.Len() > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(strings.Repeat(indent, depth))
	}

	for i, t := range toks {
		switch t.typ {
		case html.TextToken:
			switch {
			case flow > 0:
				b.WriteString(t.raw)
			case strings.TrimSpace(t.raw) != "":
				breakLine()
				b.WriteString(strings.TrimSpace(t.raw))
			}
		case html.StartTagToken:
			if flow == 0 {
				breakLine()
			}
			b.WriteString(t.raw)
			if voidElements[t.name] {
				continue
			}
			stack = append(stack, frame{t.name, t.mixed, flow > 0})
			if flow > 0 {
				continue
			}
			if t.mixed {
				flow++
			} else {
				depth++
				lastOpen = i
			}
		case html.EndTagToken:
			if len(stack) == 0 || stack[len(stack)-1].name != t.name {
				if flow == 0 {
					breakLine()
				}
				b.WriteString(t.raw)
				continue
			}
			f := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			switch {
			case f.inFlow:
			case f.mixed:
				flow--
			default:
				depth--
				if lastOpen != i-1 { // has children: close on its own line
					breakLine()
				}
			}
			b.WriteString(t.raw)
		default: // doctype, comment, self-closing tag
			if flow == 0 {
				breakLine()
			}
			b.WriteString(t.raw)
		}
	}
	b.WriteByte('\n')
	return b.Bytes()
}

// tokenize splits src and marks start tags that directly contain text.
func tokenize(src []byte) ([]token, bool) {
	z := html.NewTokenizer(bytes.NewReader(src))
	var (
		toks  []token
		stack []int // indexes into toks of open elements
	)
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			return toks, z.Err() == io.EOF
		}
		t := token{typ: tt, raw: string(z.Raw())}
		switch tt {
		case html.StartTagToken, html.EndTagToken:
			name, _ := z.TagName()
			t.name = string(name)
		}
		switch tt {
		case html.StartTagToken:
			t.mixed = verbatim[t.name]
			if !voidElements[t.name] {
				stack = append(stack, len(toks))
			}
		case html.EndTagToken:
			if n := len(stack); n > 0 && toks[stack[n-1]].name == t.name {
				stack = stack[:n-1]
			}
		case html.TextToken:
			if len(stack) > 0 && strings.TrimSpace(t.raw) != "" {
				toks[stack[len(stack)-1]].mixed = true
			}
		}
		toks = append(toks, t)
	}
}

// Middleware pretty-prints every text/html response from next.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pw := &writer{ResponseWriter: w}
		next.ServeHTTP(pw, r)
		pw.finish()
	})
}

type writer struct {
	http.ResponseWriter
	buf    bytes.Buffer
	status int
	html   bool
	sent   bool
}

func (w *writer) WriteHeader(code int) {
	if w.sent {
		return
	}
	w.sent = true
	w.status = code
	w.html = strings.HasPrefix(w.Header().Get("Content-Type"), "text/html")
	if !w.html {
		w.ResponseWriter.WriteHeader(code)
	}
}

func (w *writer) Write(p []byte) (int, error) {
	if !w.sent {
		w.WriteHeader(http.StatusOK)
	}
	if !w.html {
		return w.ResponseWriter.Write(p)
	}
	return w.buf.Write(p)
}

func (w *writer) finish() {
	if !w.html {
		return
	}
	out := Format(w.buf.Bytes())
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(w.status)
	w.ResponseWriter.Write(out)
}
