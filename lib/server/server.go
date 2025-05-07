package server

import (
	"fmt"
	"huma-app/lib/config"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	chim "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/go-chi/httplog/v2"
)

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(h.staticPath, r.URL.Path)
	fi, err := os.Stat(path)
	if os.IsNotExist(err) || fi.IsDir() {
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}
func NewSpaHandler() spaHandler {
	staticPath, indexPath := filepath.Split(config.Get().Frontend.Path)
	return spaHandler{staticPath, indexPath}
}

func setLogger() *httplog.Logger {
	logger := httplog.NewLogger(config.Get().Api.Name, httplog.Options{
		LogLevel: slog.LevelDebug,
		// JSON:             true,
		Concise:          true,
		RequestHeaders:   true,
		ResponseHeaders:  true,
		MessageFieldName: "message",
		LevelFieldName:   "severity",
		TimeFieldFormat:  time.RFC3339,
		Tags: map[string]string{
			"version": config.Get().Api.Version,
			"env":     "dev",
		},
		QuietDownRoutes: []string{
			"/",
			"/ping",
			"/openapi.yaml",
		},
		QuietDownPeriod: 10 * time.Second,
		// SourceFieldName: "source",
	})
	return logger
}

func setApiDocs(mux *chi.Mux) {
	docs := fmt.Sprintf(`<!doctype html>
    <html>
      <head>
        <title>%s %s</title>
        <meta charset="utf-8" />
        <meta
          name="viewport"
          content="width=device-width, initial-scale=1" />
      </head>
      <body>
        <script
          id="api-reference"
          data-url="/openapi.json"></script>
        <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
      </body>
    </html>`, config.Get().Name, config.Get().Version)
	mux.Get("/api/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(docs))
	})
}

func NewMux() *chi.Mux {
	logger := setLogger()
	mux := chi.NewMux()

	mux.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	mux.Use(chim.Heartbeat("/api/ping"))
	mux.Use(httplog.RequestLogger(logger))
	mux.Use(chim.Recoverer)
	setApiDocs(mux)
	mux.Handle("/*", NewSpaHandler())
	return mux
}
