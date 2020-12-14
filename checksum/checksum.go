package main

import (
	v "checksum/version"

	"crypto/md5"
	"fmt"
	"io"
	"os"

	"github.com/SebastiaanPasterkamp/gobernate"
	ge "github.com/SebastiaanPasterkamp/gonyexpress"
	pl "github.com/SebastiaanPasterkamp/gonyexpress/payload"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Checksum")

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

	c := ge.NewConsumer(rabbitmq, "checksum", 4, checksumFile)

	err := c.Run()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Shutdown()

	g.Ready()
	<-shutdown
}

func checksumFile(
	_ string, _ pl.MetaData, _ pl.Arguments, docs pl.Documents,
) (*pl.Documents, *pl.MetaData, error) {
	fh, ok := docs["input"]
	if !ok {
		return nil, nil, fmt.Errorf("did not find 'input' document")
	}

	hash := md5.New()
	if _, err := io.Copy(hash, fh.Reader()); err != nil {
		return nil, nil, err
	}

	checksum := pl.NewDocument(fmt.Sprintf("%x", hash.Sum(nil)), "text/plain", "")

	return &pl.Documents{"md5": checksum}, nil, nil
}
