package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
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
	s, err := svr.NewServer(&svr.Options{Port: 4222, JetStream: true})
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
	nc.Subscribe("math.sum", func(m *nats.Msg) {
		replyWithSum(m)
	})
}

func replyWithSum(m *nats.Msg) {
	a, b, err := parseSumRequest(m)
	if err != nil {
		resp := fmt.Sprintf("sorry, I could not understand your request: %v", err)
		m.Respond([]byte(resp))
	}

	m.Respond(sum(a, b))
}

func parseSumRequest(m *nats.Msg) (a, b int, err error) {
	str := string(m.Data)
	ss := strings.Split(str, ",")

	if len(ss) != 2 {
		err = fmt.Errorf("expected input format n1,n2: got %s", str)
		return
	}

	a, err = strconv.Atoi(ss[0])
	if err != nil {
		return
	}

	b, err = strconv.Atoi(ss[1])
	if err != nil {
		return
	}

	return
}

func sum(a, b int) []byte {
	s := fmt.Sprintf("The sum of %d and %d is %d", a, b, a+b)
	return []byte(s)
}

func requestMathSvc() {
	go func() {
		for {
			var inp string
			fmt.Println("\nEnter two numbers to be summed. eg. 2,3:")
			fmt.Scanln(&inp)

			timeout := time.Second
			m, err := nc.Request("math.sum", []byte(inp), timeout)
			if err != nil {
				log.Println(err)
			}

			fmt.Printf("%s\n", m.Data)
		}
	}()
}
