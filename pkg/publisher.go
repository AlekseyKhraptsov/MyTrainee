package main

import (
	"github.com/nats-io/nats.go"
	"log"
	"os"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	model, err := os.ReadFile("model1.json")
	if err != nil {
		log.Fatal(err)
	}

	nc.Publish("orders", model)

	defer nc.Close()

}
