package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/sealtv/worldofwisdom/internal/client"
	"github.com/sealtv/worldofwisdom/internal/pow"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	cli := client.New("localhost:8080")

	if err := cli.Connect(); err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	data, err := cli.Receive()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(data)

	challenge := strings.TrimSpace(string(data))

	// Solve the PoW challenge and find the valid nonce
	nonce, err := pow.ProofOfWork(ctx, challenge, pow.POW_DIFFICULTY)
	if err != nil {
		log.Fatal(err)
	}

	if err := cli.Send(nonce); err != nil {
		log.Fatal(err)
	}

	data, err = cli.Receive()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Result: %s", data)

}
