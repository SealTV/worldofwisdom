package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/pkg/errors"
)

type Service interface {
}

type Server struct {
	service Service
}

func New(service Service) *Server {
	return &Server{
		service: service,
	}
}

// Run starts the server and blocks until the context is canceled.
func (s *Server) Run(ctx context.Context, port int) error {
	// create a listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return errors.Wrap(err, "failed to listen")
	}
	defer listener.Close()

	conns := make(chan net.Conn)
	errs := make(chan error)

	// start accepting connections
	go s.handleIncomingConnections(conns, errs, listener)

	// handle incoming connections
	for {
		select {
		case con := <-conns:
			go s.connectionHandler(con)
		case err := <-errs:
			return err
		case <-ctx.Done():
			return nil
		}
	}
}

// handleIncomingConnections accepts incoming connections and sends them to the conns channel.
func (s *Server) handleIncomingConnections(conns chan<- net.Conn, errs chan<- error, listener net.Listener) {
	for {
		con, err := listener.Accept()
		if err != nil {
			errs <- errors.Wrap(err, "failed to accept connection")
			continue
		}

		log.Println("accepted new connection")
		conns <- con
	}
}

func (s *Server) connectionHandler(con net.Conn) {
	defer con.Close()

	// read the request
	// validate the request
	// solve the proof of work
	// send the response
	// close the connection

	fmt.Fprintln(con, "hello world")
	for {
		buf := make([]byte, 1024)

		n, err := con.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("connection closed by client")
				return
			}

			log.Printf("failed to read from connection: %v", err)
			return
		}

		log.Printf("received message: %s", string(buf[:n]))
	}
}
