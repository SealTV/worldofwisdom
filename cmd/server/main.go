package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/sealtv/worldofwisdom/internal/server"
)

func main() {
	// create a context that is canceled when the process receives an interrupt signal
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	// create a new server
	srv := server.New(nil)

	log.Println("starting server")

	// run the server
	if err := srv.Run(ctx, 8080); err != nil {
		log.Fatal(err)
	}

	log.Println("server stopped")
}
