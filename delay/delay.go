package main

import (
	v "delay/version"

	"fmt"
	"os"
	"time"

	"github.com/SebastiaanPasterkamp/gobernate"
	ge "github.com/SebastiaanPasterkamp/gonyexpress"
	pl "github.com/SebastiaanPasterkamp/gonyexpress/payload"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Delay")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set.")
	}
	rabbitmq := os.Getenv("RABBITMQ")
	if rabbitmq == "" {
		log.Fatal("RABBITMQ is not set.")
	}

	g := gobernate.New(port, v.Name, v.Release, v.Commit, v.BuildTime)

	shutdown := g.Launch()
	defer g.Shutdown()

	c := ge.NewConsumer(rabbitmq, "delay", 4, delayMessage)

	err := c.Run()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Shutdown()

	g.Ready()
	<-shutdown
}

func delayMessage(
	_ string, _ pl.MetaData, args pl.Arguments, _ pl.Documents,
) (*pl.Documents, *pl.MetaData, error) {
	duration, ok := args["duration"]
	if !ok {
		return nil, nil, fmt.Errorf("missing 'duration' in arguments")
	}

	sleep, err := time.ParseDuration(duration.(string))
	if err != nil {
		return nil, nil, err
	}

	time.Sleep(sleep)

	return nil, nil, nil
}
