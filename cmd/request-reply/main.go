package main

import (
	"fmt"
	"log"
	"time"

	svr "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/siuyin/dflt"
)

func main() {
	fmt.Println("request reply demonstration")

	embedNATS()
	mathSvcStart()
	requestMathSvc()

	select {} // wait forever
}

func embedNATS() {
	s, err := svr.NewServer(&svr.Options{Port: 4222})
	if err != nil {
		log.Fatal(err)
	}

	s.Start()
	connectSvr(s)
	log.Println("NATS server started")
}

var (
	nc  *nats.Conn
	err error
)

func connectSvr(s *svr.Server) {
	waitForSvrRdy(s)

	url := dflt.EnvString("NATS_SVR", "nats://localhost:4222")
	nc, err = nats.Connect(url)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected to NATS server: ", url)
}

func waitForSvrRdy(s *svr.Server) {
	for rdy := s.ReadyForConnections(time.Second); !rdy; rdy = s.ReadyForConnections(time.Second) {
	}
}

func mathSvcStart() {
}

func requestMathSvc() {
}
