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
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	cli := client.New(viper.GetString("server"))

	if err := cli.Connect(); err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	challenge, err := cli.Receive()
	if err != nil {
		log.Fatalf("receive challenge err: %v", err)
	}

	pow := pow.NewPoW(viper.GetInt("pow_difficulty"))

	// Solve the PoW challenge and find the valid nonce
	nonce, err := pow.ProofOfWork(ctx, challenge)
	if err != nil {
		log.Fatalf("failed to solve PoW challenge: %v", err)
	}

	if err := cli.Send(nonce); err != nil {
		log.Fatalf("Send problem solving err: %v", err)
	}

	quoute, err := cli.Receive()
	if err != nil {
		log.Fatalf("receive result err: %v", err)
	}

	log.Printf("Result: '%s'", quoute)

}
