package main

import (
	"context"
	"fmt"
	"huma-app/controller"
	"huma-app/lib/middleware"
	"huma-app/lib/security"
	"huma-app/store"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	chim "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

const staticPath = "frontend/build"
const indexPath = "index.html"

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

//go:generate sqlc generate

// Options for the CLI.
type Options struct {
	JwtKey      string `doc:"Secret key for jwt generation" name:"jwtkey"`
	Name        string `doc:"App name" name:"name" default:"My API"`
	Version     string `doc:"App version" default:"1.0.0"`
	Port        int    `doc:"Port to listen on" short:"p" default:"8888"`
	Debug       bool   `doc:"Enable debug" short:"d" default:"false"`
	StoragePath string `doc:"Path to Sqlite DB" name:"path" default:"./storage.db"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err.Error())
	}

	cli := humacli.New(func(hooks humacli.Hooks, opts *Options) {

		logger := setupLogger(opts)

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

		mux.Use(chim.Heartbeat("/ping"))
		mux.Use(httplog.RequestLogger(logger))

		db := store.InitDB(opts.StoragePath)
		store := store.New(db)
		security := security.NewSecurity(opts.JwtKey)
		users := controller.NewUserResource(store, security)

		config := huma.DefaultConfig(opts.Name, opts.Version)
		config.DocsPath = "/api/docs"
		/* config.CreateHooks = []func(huma.Config) huma.Config{
		func(c huma.Config) huma.Config {
			return c
		}} */
		config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
			"Bearer": {
				Type:         "http",
				Scheme:       "bearer",
				BearerFormat: "JWT",
				In:           "cookie",
			},
		}

		api = humachi.New(mux, config)
		api.UseMiddleware(middleware.JwtAuthMiddleware(api, security))

		huma.AutoRegister(api, users)
		spa := spaHandler{staticPath: staticPath, indexPath: indexPath}
		mux.Handle("/*", spa)

		server := http.Server{
			Addr:    fmt.Sprintf(":%d", opts.Port),
			Handler: mux,
		}

		slog.Info("Server started on ", "localhost:", opts.Port)

		hooks.OnStart(func() {
			server.ListenAndServe()
		})

		hooks.OnStop(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			server.Shutdown(ctx)
			db.Close()
		})

	})
	cli.Root().AddCommand(&cobra.Command{
		Use:   "openapi",
		Short: "Print the OpenAPI spec",
		Run: func(cmd *cobra.Command, args []string) {
			b, err := api.OpenAPI().MarshalJSON()
			if err != nil {
				panic(err)
			}
			path := "openapi.json"
			f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

			if err != nil {
				panic(err)
			}
			_, err = f.Write(b)
			if err != nil {
				panic(err)
			}

			err = f.Close()
			if err != nil {
				panic(err)
			}

			fmt.Printf("spec saved to %s\n", path)
		},
	})
	cli.Run()

}

func setupLogger(opts *Options) *httplog.Logger {
	logger := httplog.NewLogger(opts.Name, httplog.Options{
		LogLevel: slog.LevelDebug,
		// JSON:             true,
		Concise:          true,
		RequestHeaders:   true,
		ResponseHeaders:  true,
		MessageFieldName: "message",
		LevelFieldName:   "severity",
		TimeFieldFormat:  time.RFC3339,
		Tags: map[string]string{
			"version": opts.Version,
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
