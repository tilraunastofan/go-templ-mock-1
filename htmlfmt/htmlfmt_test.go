package htmlfmt

import "testing"

func TestFormat(t *testing.T) {
	tests := []struct{ name, in, want string }{
		{
			"nesting",
			`<div><p class="a">Hi</p><ul><li>one</li><li>two</li></ul></div>`,
			"<div>\n  <p class=\"a\">Hi</p>\n  <ul>\n    <li>one</li>\n    <li>two</li>\n  </ul>\n</div>\n",
		},
		{
			"inline stays inline inside text",
			`<p>Hello <b>big</b> world</p>`,
			"<p>Hello <b>big</b> world</p>\n",
		},
		{
			"void and empty elements",
			`<div><input type="text"><span></span></div>`,
			"<div>\n  <input type=\"text\">\n  <span></span>\n</div>\n",
		},
		{
			"verbatim",
			"<div><pre>a\n  <b>b</b></pre><script>if (a<b) {x()}</script></div>",
			"<div>\n  <pre>a\n  <b>b</b></pre>\n  <script>if (a<b) {x()}</script>\n</div>\n",
		},
		{
			"entities kept",
			`<p>a &amp; b</p>`,
			"<p>a &amp; b</p>\n",
		},
		{
			"doctype",
			`<!DOCTYPE html><html><body><p>x</p></body></html>`,
			"<!DOCTYPE html>\n<html>\n  <body>\n    <p>x</p>\n  </body>\n</html>\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := string(Format([]byte(tt.in))); got != tt.want {
				t.Errorf("got:\n%q\nwant:\n%q", got, tt.want)
			}
		})
	}
}
