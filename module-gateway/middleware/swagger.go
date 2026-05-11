package middleware

import (
	"context"
	"html/template"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	swaggerFiles "github.com/swaggo/files/v2"
	swag "github.com/swaggo/swag/v2"
)

// SwaggerConfig holds the Swagger UI configuration.
type SwaggerConfig struct {
	// URL points to the API definition JSON. Default: "doc.json".
	URL          string
	Title        string
	DocExpansion string
	InstanceName string
	DeepLinking  bool
}

// SwaggerHandler returns a Hertz HandlerFunc that serves the Swagger UI.
// Mount it on a wildcard route, e.g.:
//
//	h.GET("/swagger/*any", middleware.SwaggerHandler())
func SwaggerHandler(opts ...func(*SwaggerConfig)) app.HandlerFunc {
	cfg := &SwaggerConfig{
		URL:          "doc.json",
		Title:        "Swagger UI",
		DocExpansion: "list",
		InstanceName: swag.Name,
		DeepLinking:  true,
	}
	for _, o := range opts {
		o(cfg)
	}

	indexTmpl, _ := template.New("swagger_index.html").Parse(swaggerIndexTpl)

	return func(c context.Context, ctx *app.RequestContext) {
		if string(ctx.Request.Method()) != consts.MethodGet {
			ctx.AbortWithStatus(http.StatusMethodNotAllowed)
			return
		}

		// Extract the part after the prefix (e.g. /swagger/)
		rawPath := string(ctx.Request.URI().Path())
		// Find the last occurrence of "swagger/" and extract what comes after
		idx := strings.LastIndex(rawPath, "swagger/")
		var filePath string
		if idx >= 0 {
			filePath = rawPath[idx+len("swagger/"):]
		}
		if filePath == "" {
			filePath = "index.html"
		}

		switch path.Ext(filePath) {
		case ".html":
			ctx.Header("Content-Type", "text/html; charset=utf-8")
		case ".css":
			ctx.Header("Content-Type", "text/css; charset=utf-8")
		case ".js":
			ctx.Header("Content-Type", "application/javascript")
		case ".png":
			ctx.Header("Content-Type", "image/png")
		case ".json":
			ctx.Header("Content-Type", "application/json; charset=utf-8")
		}

		switch filePath {
		case "index.html":
			_ = indexTmpl.Execute(ctx, cfg)
		case "doc.json":
			doc, err := swag.ReadDoc(cfg.InstanceName)
			if err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			_, _ = ctx.Write([]byte(doc))
		default:
			data, err := fs.ReadFile(swaggerFiles.FS, filePath)
			if err != nil {
				ctx.AbortWithStatus(http.StatusNotFound)
				return
			}
			_, _ = ctx.Write(data)
		}
	}
}

// Option helpers ──────────────────────────────────────────────────────────────

func SwaggerURL(url string) func(*SwaggerConfig) {
	return func(c *SwaggerConfig) { c.URL = url }
}

func SwaggerTitle(title string) func(*SwaggerConfig) {
	return func(c *SwaggerConfig) { c.Title = title }
}

func SwaggerInstanceName(name string) func(*SwaggerConfig) {
	return func(c *SwaggerConfig) { c.InstanceName = name }
}

// ── HTML template ─────────────────────────────────────────────────────────────

const swaggerIndexTpl = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>{{.Title}}</title>
  <link rel="stylesheet" type="text/css" href="swagger-ui.css">
  <link rel="icon" type="image/png" href="favicon-32x32.png" sizes="32x32"/>
  <link rel="icon" type="image/png" href="favicon-16x16.png" sizes="16x16"/>
</head>
<body>
<div id="swagger-ui"></div>
<script src="swagger-ui-bundle.js"></script>
<script src="swagger-ui-standalone-preset.js"></script>
<script>
window.onload = function() {
  const ui = SwaggerUIBundle({
    url: "{{.URL}}",
    dom_id: '#swagger-ui',
    deepLinking: {{.DeepLinking}},
    docExpansion: "{{.DocExpansion}}",
    presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
    plugins: [SwaggerUIBundle.plugins.DownloadUrl],
    layout: "StandaloneLayout"
  });
  window.ui = ui;
};
</script>
</body>
</html>`

