package main

import (
	"math/rand"
	v "unstable/version"

	"fmt"
	"os"

	"github.com/SebastiaanPasterkamp/gobernate"
	ge "github.com/SebastiaanPasterkamp/gonyexpress"
	pl "github.com/SebastiaanPasterkamp/gonyexpress/payload"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Unstable")

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

	c := ge.NewConsumer(rabbitmq, "unstable", 4, randomlyFail)

	err := c.Run()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Shutdown()

	g.Ready()
	<-shutdown
}

func randomlyFail(
	_ string, _ pl.MetaData, args pl.Arguments, _ pl.Documents,
) (*pl.Documents, *pl.MetaData, error) {
	rate, ok := args["error_rate"]
	if !ok {
		return nil, nil, fmt.Errorf("missing 'error_rate' in arguments")
	}

	if rand.Float64() <= rate.(float64) {
		return nil, nil, fmt.Errorf("randomly failed for no reason")
	}

	return nil, nil, nil
}
