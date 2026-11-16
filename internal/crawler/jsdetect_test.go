package crawler

import (
	"strings"
	"testing"

	"github.com/dillonlara115/barracudaseo/pkg/models"
)

func detect(t *testing.T, html string) *models.JSInfo {
	t.Helper()
	doc := docOf(t, html)
	return detectJS(doc, ExtractContent(doc))
}

func TestDetectJS(t *testing.T) {
	article := strings.Repeat("Real server rendered words. ", 40)

	tests := []struct {
		name           string
		html           string
		wantNil        bool
		framework      string
		clientRendered bool
	}{
		{
			name:    "plain static page",
			html:    `<html><body><main><p>` + article + `</p></main><script src="/analytics.js"></script></body></html>`,
			wantNil: true,
		},
		{
			name: "client-rendered React shell",
			html: `<html><body><noscript>You need to enable JavaScript to run this app.</noscript>
				<div id="root"></div><script src="/static/js/main.abc123.js"></script></body></html>`,
			framework:      "React",
			clientRendered: true,
		},
		{
			name: "server-rendered Next.js page",
			html: `<html><body><div id="__next"><main><p>` + article + `</p></main></div>
				<script id="__NEXT_DATA__" type="application/json">{}</script>
				<script src="/_next/static/chunks/main.js"></script></body></html>`,
			framework:      "Next.js",
			clientRendered: false,
		},
		{
			name: "Angular shell",
			html: `<html><body><app-root></app-root><script src="/runtime.js"></script>
				<script src="/main.js"></script></body></html>`,
			framework:      "Angular",
			clientRendered: true,
		},
		{
			name: "empty body with several scripts and no known mount point",
			html: `<html><body><div class="container"></div><script src="/a.js"></script>
				<script src="/b.js"></script><script src="/c.js"></script></body></html>`,
			clientRendered: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := detect(t, tt.html)
			if tt.wantNil {
				if info != nil {
					t.Fatalf("info = %+v, want nil", info)
				}
				return
			}
			if info == nil {
				t.Fatal("info = nil, want a result")
			}
			if info.Framework != tt.framework {
				t.Errorf("Framework = %q, want %q", info.Framework, tt.framework)
			}
			if info.ClientRendered != tt.clientRendered {
				t.Errorf("ClientRendered = %v, want %v (signals: %v)", info.ClientRendered, tt.clientRendered, info.Signals)
			}
		})
	}
}
