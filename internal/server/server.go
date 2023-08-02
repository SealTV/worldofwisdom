package server

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/pkg/errors"
	"github.com/sealtv/worldofwisdom/internal/app"
)

type Service interface {
	ProcessClient(c app.Clienter) error
}

type Server struct {
	listener net.Listener
	service  Service
}

func New(listener net.Listener, service Service) *Server {
	return &Server{
		listener: listener,
		service:  service,
	}
}

// Run starts the server and blocks until the context is canceled.
func (s *Server) Run(ctx context.Context) error {
	conns := make(chan net.Conn)
	errs := make(chan error)

	// start accepting connections
	go handleIncomingConnections(s.listener, conns, errs)

	// wait for all connections to finish
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	// handle incoming connections
	for {
		select {
		case conn := <-conns:
			wg.Add(1)

			// handle the connection
			go func(conn net.Conn) {
				defer wg.Done()
				defer conn.Close()

				cli := NewClient(conn)

				// Handle the client
				if err := s.service.ProcessClient(cli); err != nil {
					log.Printf("failed to handle client: %v", err)
					return
				}
			}(conn)
		case err := <-errs:
			return err
		case <-ctx.Done():
			return nil
		}
	}
}

// handleIncomingConnections accepts incoming connections and sends them to the conns channel.
func handleIncomingConnections(listener net.Listener, conns chan<- net.Conn, errs chan<- error) {
	defer close(conns)
	defer close(errs)

	for {
		con, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Println("listener closed")

				errs <- nil
				return
			}

			errs <- errors.Wrap(err, "failed to accept connection")
			continue
		}

		log.Println("accepted new connection")
		conns <- con
	}
}
