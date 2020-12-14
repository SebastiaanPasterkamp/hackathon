package main

import (
	v "attach/version"

	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/SebastiaanPasterkamp/gobernate"
	ge "github.com/SebastiaanPasterkamp/gonyexpress"
	pl "github.com/SebastiaanPasterkamp/gonyexpress/payload"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Attach")

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

	c := ge.NewConsumer(rabbitmq, "attach", 4, attachFile)

	err := c.Run()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Shutdown()

	g.Ready()
	<-shutdown
}

func attachFile(
	_ string, _ pl.MetaData, _ pl.Arguments, docs pl.Documents,
) (*pl.Documents, *pl.MetaData, error) {
	name, ok := docs["filename"]
	if !ok {
		return nil, nil, fmt.Errorf("did not find 'filename' document")
	}

	path := filepath.Join("data", filepath.Base(name.Data))
	fh, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer fh.Close()

	d := pl.InitDocument("text/plain", pl.Base64Encoding)
	dw := d.WriteCloser()
	defer dw.Close()

	if _, err := io.Copy(dw, fh); err != nil {
		return nil, nil, err
	}

	return &pl.Documents{"input": d}, nil, nil
}
