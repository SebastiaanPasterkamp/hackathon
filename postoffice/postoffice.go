package main

import (
	v "postoffice/version"

	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/SebastiaanPasterkamp/gobernate"
	ge "github.com/SebastiaanPasterkamp/gonyexpress"
	pl "github.com/SebastiaanPasterkamp/gonyexpress/payload"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func main() {
	log.Info("Post Office")

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

	p := ge.NewProducer(rabbitmq, "postoffice")

	msgs, err := p.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close()

	var wg = sync.WaitGroup{}

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go handleMessages(&p, msgs, &wg, shutdown)
	}

	g.Ready()

	<-shutdown
	wg.Wait()
}

func handleMessages(p *ge.Component, msgs <-chan amqp.Delivery, wg *sync.WaitGroup, shutdown chan bool) {
	defer wg.Done()

	var running = true
	for running {
		select {
		case <-shutdown:
			running = false
			continue

		case d := <-msgs:
			msg, err := pl.MessageFromByteSlice(d.Body)

			if err != nil {
				log.Errorf("%s - Bad message: %+v in %+v\n",
					d.CorrelationId, err, d.Body)
				d.Nack(false, false)
				continue
			}

			r, err := readRoute(msg.Routing.Name)
			if err != nil {
				log.Errorf("%s - Unknown route %q: %+v\n",
					d.CorrelationId, msg.Routing.Name, err)
				d.Nack(false, false)
				continue
			}
			if d.ReplyTo != "" {
				r.Slip = append(r.Slip, pl.Step{Queue: d.ReplyTo})
			}

			msg.Routing = *r
			p.SendMessage(*msg)

			d.Ack(false)
		}
	}
}

func readRoute(name string) (*pl.Routing, error) {
	path := filepath.Join("data", filepath.Base(name)+".json")
	fh, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	data, err := ioutil.ReadAll(fh)
	if err != nil {
		return nil, err
	}

	var r pl.Routing
	err = json.Unmarshal(data, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}
