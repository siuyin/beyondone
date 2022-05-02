package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

var lst = []string{
	"marraige proposal accepted",
	"marraige registration witnesses invited",
	"marraige registered",
	"official photographer hired",
	"wedding gown for bride selected",
	"wedding dinner guest list finalised",
	"wedding dinner guests invited",
	"wedding dinner menus finalised",
	"wedding cake selected",
	"wedding dinner restaurant booked",
	"seating plan finalised",
}

var (
	nc  *nats.Conn
	js  nats.JetStreamContext
	err error
)

func main() {
	fmt.Println("temporal decoupling with event streaming")

	embedNATS()
	defineEventStream(lst)

	select {} // wait forever
}

func defineEventStream(lst []string) {
	if !initStream() {
		return
	}

	addStream()
	feedEvents(lst)

	log.Println("event stream initialised")
}

func initStream() bool {
	if len(os.Args) != 2 {
		return false
	}

	if os.Args[1] != "start" {
		return false
	}

	return true
}

func addStream() {

	if _, err := js.AddStream(&nats.StreamConfig{
		Name: "wedding", Subjects: []string{"wedding", "wedding.>"},
	}); err != nil {
		log.Fatal(err)
	}
}

func feedEvents(lst []string) {
	for _, item := range lst {
		js.PublishAsync("wedding.list", []byte(item))
	}
}
