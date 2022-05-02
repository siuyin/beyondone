package main

import (
	"log"
	"time"

	svr "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/siuyin/dflt"
)

func embedNATS() {
	s, err := svr.NewServer(&svr.Options{Port: 4222, JetStream: true})
	if err != nil {
		log.Fatal(err)
	}

	s.Start()
	connectSvr(s)
	initJetStream()
	log.Println("NATS server started")
}

func connectSvr(s *svr.Server) {
	waitForSvrRdy(s)

	url := dflt.EnvString("NATS_SVR", "nats://localhost:4222")
	nc, err = nats.Connect(url)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected to NATS server: ", url)
}

func initJetStream() {
	js, err = nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("JetStream event streaming persistence initialised")
}

func waitForSvrRdy(s *svr.Server) {
	for rdy := s.ReadyForConnections(time.Second); !rdy; rdy = s.ReadyForConnections(time.Second) {
	}
}
