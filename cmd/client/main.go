package main

import (
	"log"

	"github.com/sealtv/worldofwisdom/internal/client"
)

func main() {
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

	cli.Send("hello world from client")
}
