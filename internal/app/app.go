package app

import (
	"context"
	"log"
	"time"
)

type Clienter interface {
	ReadWithTimeout(ctx context.Context, timeout time.Duration) (string, error)
	Read(ctx context.Context) (string, error)
	Write(msg string) error
}

type PoWer interface {
	GetChallenge() string
	IsValid(input string) bool
}

type WisdomBooker interface {
	GetRandomQuote() string
}

type App struct {
	timeout time.Duration
	power   PoWer
	wb      WisdomBooker
}

func NewApp(power PoWer, wb WisdomBooker) *App {
	return &App{
		timeout: 5 * time.Second,
		power:   power,
		wb:      wb,
	}
}

func (a *App) ProcessClient(cli Clienter) error {
	challenge := a.power.GetChallenge()
	if err := cli.Write(challenge); err != nil {
		log.Printf("cannot send challenge: %v", err)
		return nil
	}

	clientResponse, err := cli.ReadWithTimeout(context.Background(), a.timeout)
	if err != nil {
		log.Printf("error on read client response: %v", err)
		return nil
	}

	// Verify the PoW response from the client
	if !a.power.IsValid(challenge + clientResponse) {
		if err := cli.Write("Invalid PoW response"); err != nil {
			log.Printf("failed to write to connection: %v", err)
		}

		return nil
	}

	if err := cli.Write(a.wb.GetRandomQuote()); err != nil {
		log.Printf("failed to write to connection: %v", err)
	}

	return nil
}
