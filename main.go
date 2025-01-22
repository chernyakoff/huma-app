package main

import (
	"context"
	"database/sql"
	"fmt"
	"huma-app/controller"
	"huma-app/lib/config"
	"huma-app/lib/middleware"
	"huma-app/lib/security"
	"huma-app/store"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alecthomas/kong"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	chim "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	"github.com/ne-sachirou/go-graceful"
	"github.com/ne-sachirou/go-graceful/gracefulhttp"
)

var api huma.API

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
	staticPath, indexPath := filepath.Split(config.Get().Server.Spa)
	return spaHandler{staticPath, indexPath}
}

//go:generate sqlc generate

func setupMux() *chi.Mux {
	logger := setupLogger()
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
	mux.Use(chim.Heartbeat("/ping"))
	mux.Use(httplog.RequestLogger(logger))
	mux.Get("/api/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(docs))
	})

	mux.Handle("/*", NewSpaHandler())
	return mux
}

func setupApi(db *sql.DB, mux *chi.Mux) huma.API {
	store := store.New(db)
	security := security.NewSecurity()
	users := controller.NewUserResource(store, security)
	config := huma.DefaultConfig(config.Get().Name, config.Get().Version)
	config.DocsPath = "" //"/api/docs"
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"Bearer": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
			In:           "cookie",
		},
	}
	config.CreateHooks = []func(huma.Config) huma.Config{
		func(c huma.Config) huma.Config { return c },
	}
	api = humachi.New(mux, config)
	api.UseMiddleware(middleware.JwtAuthMiddleware(api, security))
	huma.AutoRegister(api, users)
	return api
}

func setupLogger() *httplog.Logger {
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

func MustWriteFile(path string, content []byte) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic("Ошибка при открытии файла: " + path)
	}
	defer file.Close()
	_, err = file.Write(content)
	if err != nil {
		panic("Ошибка при записи в файл: " + path)
	}

}

type KongContext struct {
	Debug bool
}

type ServeCmd struct {
}

func (cmd *ServeCmd) Run(ctx *KongContext) error {

	db := store.InitDB()
	mux := setupMux()
	api = setupApi(db, mux)
	b, err := api.OpenAPI().MarshalJSON()
	if err != nil {
		panic(err)
	}
	specPath := "openapi.json"
	MustWriteFile(specPath, b)

	slog.Info("Openapi spec saved to " + specPath)

	addr := fmt.Sprintf("%s:%d", config.Get().Server.Host, config.Get().Server.Port)
	// defer db.Close()
	if err := gracefulhttp.ListenAndServe(
		context.Background(),
		addr,
		mux,
		graceful.GracefulShutdownTimeout(time.Second),
	); err != nil {
		panic(err)
	}

	slog.Info("Server started on " + addr)

	return nil
}

var CLI struct {
	Serve ServeCmd `cmd:"serve" default:"1" help:"Start Server"`
}

func main() {

	kong.Parse(&CLI).Run(&KongContext{})

}
