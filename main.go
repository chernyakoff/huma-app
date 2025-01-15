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
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	chim "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"github.com/joho/godotenv"
)

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
		//setupLogger(opts)
		logger := setupLogger(opts)
		mux := chi.NewMux()
		mux.Use(chim.Heartbeat("/ping"))
		mux.Use(httplog.RequestLogger(logger))

		mux.Handle("/zopa/*", http.FileServer(http.Dir("./frontend/build")))

		db := store.InitDB(opts.StoragePath)

		store := store.New(db)
		config := huma.DefaultConfig(opts.Name, opts.Version)
		config.DocsPath = "/api/docs"
		config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
			"Bearer": {
				Type:         "http",
				Scheme:       "bearer",
				BearerFormat: "JWT",
				In:           "cookie",
			},
		}
		security := security.NewSecurity(opts.JwtKey)
		api := humachi.New(mux, config)
		api.UseMiddleware(middleware.JwtAuthMiddleware(api, security))
		users := controller.NewUserResource(store, security)
		huma.AutoRegister(api, users)
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
