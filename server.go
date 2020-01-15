package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"

	"aqwari.net/net/styx"
)

type server struct {
	services []*service
	ctx      context.Context
	cfg      *config
	events   chan *event
	inputs   chan interface{}
	clients  chan interface{}
	controls chan interface{}
	sync.Mutex
	errors []error
}

type service struct {
	addr string
	name string
}

type message struct {
	service string
	data    string
	buff    string
}

type fileHandler struct {
	fn func(msg *message) (interface{}, error)
}

var handlers = make(map[string]*fileHandler)

func addFileHandler(path string, fh *fileHandler) {
	handlers[path] = fh
}

func newServer(ctx context.Context, cfg *config) (*server, error) {
	events, err := listenEvents(ctx, cfg)
	if err != nil {
		return nil, err
	}
	var services []*service
	for _, svc := range cfg.listServices() {
		service := &service{
			addr: cfg.getAddress(svc),
			name: svc,
		}
		services = append(services, service)
	}
	s := &server{
		services: services,
		events:   events,
		inputs:   make(chan interface{}),
		controls: make(chan interface{}),
		clients:  make(chan interface{}),
		ctx:      ctx,
		cfg:      cfg,
	}
	return s, nil
}

func (s *server) listenEvents() {
	for event := range s.events {
		s.Lock()
		defer s.Unlock()
		// We want to update our service/server here for when we generate content for tabs.
		// Event name is where it comes from, event lines are >0 lines of events
		fmt.Println(event.lines)
	}
}

func (s *server) start() {
	for _, svc := range s.services {
		go run(s, svc)
	}
}

func run(s *server, srv *service) {
	h := styx.HandlerFunc(func(sess *styx.Session) {
		uuid := rand.Int63()
		s.clients <- &client{
			target: srv.name,
			uuid:   uuid,
		}
		for sess.Next() {
			handleReq(s, sess.Request())
		}
	})
	t := &styx.Server{
		Addr:    srv.addr + ":564",
		Handler: h,
		//Auth: auth,
	}
	var err error
	if *usetls {
		err = t.ListenAndServeTLS(*certfile, *keyfile)
	} else {
		err = t.ListenAndServe()
	}
	if err != nil {
		log.Fatal(err)
	}
}

func handleReq(s *server, req styx.Request) {
	switch msg := req.(type) {
	case styx.Twalk:
		msg.Rwalk(os.Stat(msg.Path()))
	case styx.Topen:
	case styx.Tstat:
	}
}
