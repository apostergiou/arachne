package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/miekg/dns"
)

// Server is the component that performs the actual work (receives normal DNS
// requests and forwards them to an upstream DNS server using DNS-over-TLS).
type Server struct {
	Log *log.Logger

	srv *dns.Server
	cfg *Config
}

// NewServer accepts a non-nil configuration, an optional logger and
// returns a new Server.
// If logger is nil, server logs are disabled.
func NewServer(cfg *Config, logger *log.Logger) (*Server, error) {
	if cfg == nil {
		return nil, errors.New("config cannot be nil")
	}

	if logger == nil {
		logger = log.New(ioutil.Discard, "", 0)
	}

	s := new(Server)
	s.srv = &dns.Server{
		Addr: cfg.Listen, Net: cfg.Network, Handler: dns.HandlerFunc(s.Handler)}
	s.cfg = cfg
	s.Log = logger

	return s, nil
}

// ListenAndServe listens on the TCP network address s.Listen and handles
// requests on incoming connections.
func (s *Server) ListenAndServe() error {
	s.Log.Printf("Listening on %s (%s)", s.cfg.Listen, s.cfg.Network)
	s.Log.Printf("Forwarding to %s", s.cfg.Upstream)
	s.Log.Printf("Configuration: %#v", s.cfg)
	go func() {
		err := s.srv.ListenAndServe()
		if err != nil {
			s.Log.Fatalf("Failed to start the server: %s", err)
		}
	}()

	sgC := make(chan os.Signal)
	signal.Notify(sgC, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case sg := <-sgC:
			fmt.Printf("Signal (%d) received, gracefully exiting...\n", sg)
			os.Exit(1)
		}
	}
}

// Handler forwards the DNS query to upstream DNS server.
func (s *Server) Handler(w dns.ResponseWriter, msg *dns.Msg) {
	s.Log.Printf("Received DNS query for: %s", strings.TrimRight(msg.Question[0].Name, "."))

	c := new(dns.Client)
	c.Net = "tcp-tls"
	c.Timeout = 3 * time.Second

	s.Log.Printf("Forwarding the query to: %s", s.cfg.Upstream)
	r, _, err := c.Exchange(msg, s.cfg.Upstream)
	if err != nil {
		s.Log.Fatalf("Failed to handle the query: %s", err)
	}

	s.Log.Printf("Upstream answer: %v", r.Answer)
	w.WriteMsg(r)
}
