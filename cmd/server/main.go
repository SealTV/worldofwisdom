package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/sealtv/worldofwisdom/internal/app"
	"github.com/sealtv/worldofwisdom/internal/pow"
	"github.com/sealtv/worldofwisdom/internal/server"
	"github.com/sealtv/worldofwisdom/internal/wisdombook"
)

func main() {
	// create a context that is canceled when the process receives an interrupt signal
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	// create a listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	book := wisdombook.NewWisdomBook([]string{
		"Be yourself; everyone else is already taken.",
		"A room without books is like a body without a soul.",
		"So many books, so little time.",
		"You only live once, but if you do it right, once is enough.",
		"Be the change that you wish to see in the world.",
		"In three words I can sum up everything I've learned about life: it goes on.",
	})

	// create a new server
	srv := server.New(listener, app.NewApp(pow.NewPoW(), book))

	log.Println("starting server")
	// run the server
	if err := srv.Run(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("server stopped")
}
