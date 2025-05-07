package main

import (
	"context"
	"flag"
	"fmt"
	"huma-app/lib/api"
	"huma-app/lib/config"
	"huma-app/lib/server"
	"huma-app/store"
	"log/slog"
	"time"

	"github.com/ne-sachirou/go-graceful"
	"github.com/ne-sachirou/go-graceful/gracefulhttp"
)

func main() {

	spec := flag.Bool("spec", false, "write openapi.json")
	flag.Parse()

	db := store.InitDB()
	mux := server.NewMux()
	api.NewApi(db, mux)
	if *spec {
		api.MustWriteSpec(db, mux)
		return
	}

	addr := fmt.Sprintf("%s:%d",
		config.Get().Server.Host,
		config.Get().Server.Port,
	)
	defer db.Close()
	slog.Info("Server started on " + addr)

	if err := gracefulhttp.ListenAndServe(
		context.Background(),
		addr,
		mux,
		graceful.GracefulShutdownTimeout(time.Second),
	); err != nil {
		panic(err)
	}

}
