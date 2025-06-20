package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"template/pgk"
)

const configPath = "configs/config.yml"

func main() {
	// flags in future if needed
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app := pgk.NewApp(configPath, ctx)

	if err := app.Start(ctx); err != nil {
		log.Fatal(err)
	}
}
