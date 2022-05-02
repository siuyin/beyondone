package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

var lst = []string{
	"01 marraige proposal accepted",
	"02 marraige registration witnesses invited",
	"03 marraige registered",
	"04 official photographer hired",
	"05 wedding gown for bride selected",
	"06 wedding dinner guest list finalised",
	"07 wedding dinner guests invited",
	"08 wedding dinner menus finalised",
	"09 wedding cake selected",
	"10 wedding dinner restaurant booked",
	"11 seating plan finalised",
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

	procEvtStream() // uncomment me

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
	if _, err := js.AddStream(
		&nats.StreamConfig{
			Name:     "wedding",
			Subjects: []string{"wedding", "wedding.>"},
		},
	); err != nil {
		log.Fatal(err)
	}
}

func feedEvents(lst []string) {
	for _, item := range lst {
		js.PublishAsync("wedding.list", []byte(item))
	}
}

func procEvtStream() {
	if !weddingStreamPresent() {
		return
	}

	go func() {
		for {
			msg, err := fetch()
			if err != nil {
				continue
			}

			fmt.Println("weddingPlanner: ", string(msg.Data))
			msg.Ack()
			time.Sleep(5 * time.Second) // even machine wedding planners take time
		}
		fmt.Println("forever loop exited. This should not happen")
	}()
}

func weddingStreamPresent() bool {
	_, err := js.StreamInfo("wedding")
	if err != nil {
		return false
	}
	return true
}

func fetch() (m *nats.Msg, err error) {
	sub, err := js.PullSubscribe("wedding.list", "weddingPlanner")
	if err != nil {
		log.Println(err)
		return
	}

	msgs, err := sub.Fetch(1)
	if err != nil && err.Error() == "nats: timeout" {
		return
	}
	if err != nil {
		log.Println(err)
		return
	}

	if len(msgs) != 1 {
		return m, fmt.Errorf("unexpected number of messages returned: %n", len(msgs))
	}

	return msgs[0], nil
}
