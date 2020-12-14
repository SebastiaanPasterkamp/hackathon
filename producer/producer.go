package main

import (
	"flag"
	"math/rand"
	"os"

	ge "github.com/SebastiaanPasterkamp/gonyexpress"
	pl "github.com/SebastiaanPasterkamp/gonyexpress/payload"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Print("Producer")

	total := flag.Int("total", 1, "Number of messages to send.")
	flag.Parse()

	rabbitmq := os.Getenv("RABBITMQ")
	if rabbitmq == "" {
		log.Fatal("RABBITMQ is not set.")
	}

	p := ge.NewProducer(rabbitmq, "")

	_, err := p.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close()

	routes := []string{
		"fast-and-easy",
		"fast-but-risky",
		"slow-and-risky",
		"slow-but-steady",
	}
	filenames := []string{
		"first.txt",
		"second.txt",
		"third.txt",
	}

	for i := 0; i < *total; i++ {
		msg := pl.NewMessage(
			pl.Routing{
				Name:     routes[rand.Intn(len(routes))],
				Position: 0,
				Slip: []pl.Step{
					{Queue: "postoffice"},
				},
			},
			pl.MetaData{
				"origin": "producer",
			},
			pl.Documents{
				"filename": pl.NewDocument(
					filenames[rand.Intn(len(filenames))],
					"text/plain",
					"",
				),
			},
		)

		log.Infof("%+v", msg)

		err := p.SendMessage(msg)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s: Sent %+v\n", msg.TraceID, msg)
	}

	log.Println("Successfully Published Message(s) to Queue")
}
