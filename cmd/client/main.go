package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/sealtv/worldofwisdom/internal/client"
	"github.com/sealtv/worldofwisdom/internal/pow"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.String("server", "localhost:8080", "port to listen on")
	pflag.Int("pow_difficulty", 5, "PoW difficulty")

	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()
}

func main() {
	// create a context that is canceled when the process receives an interrupt signal
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	// create a new client
	cli := client.New(viper.GetString("server"))

	if err := cli.Connect(); err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// Receive the PoW challenge
	challenge, err := cli.Receive()
	if err != nil {
		log.Fatalf("receive challenge err: %v", err)
	}

	// Create a new PoW solver
	pow := pow.NewPoW(viper.GetInt("pow_difficulty"))

	// Solve the PoW challenge and find the valid nonce
	nonce, err := pow.ProofOfWork(ctx, challenge)
	if err != nil {
		log.Fatalf("failed to solve PoW challenge: %v", err)
	}

	// Send the nonce to the server
	if err := cli.Send(nonce); err != nil {
		log.Fatalf("Send problem solving err: %v", err)
	}

	// Receive the result
	result, err := cli.Receive()
	if err != nil {
		log.Fatalf("receive result err: %v", err)
	}

	log.Printf("Result: '%s'", result)

}
