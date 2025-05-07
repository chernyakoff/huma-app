package api

import (
	"database/sql"
	"huma-app/lib/config"
	"huma-app/lib/handlers"
	"huma-app/lib/middleware"
	"huma-app/lib/security"

	"huma-app/store"
	"log/slog"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func NewApi(db *sql.DB, mux *chi.Mux) huma.API {
	store := store.New(db)
	security := security.NewSecurity()
	handlers := handlers.NewApiHandlers(store, security)
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
	api := humachi.New(mux, config)
	api.UseMiddleware(middleware.RealIpMiddleware)
	api.UseMiddleware(middleware.JwtAuthMiddleware(api, security))
	huma.AutoRegister(api, handlers)

	return api
}

func MustWriteSpec(db *sql.DB, mux *chi.Mux) {
	api := NewApi(db, mux)
	b, err := api.OpenAPI().MarshalJSON()
	if err != nil {
		panic(err)
	}
	path := "openapi.json"

	slog.Info("Openapi spec saved to " + path)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic("Ошибка при открытии файла: " + path)
	}
	defer file.Close()
	_, err = file.Write(b)
	if err != nil {
		panic("Ошибка при записи в файл: " + path)
	}

}
