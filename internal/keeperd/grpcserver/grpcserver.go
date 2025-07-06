// Package grpcserver implements handy wrap around gRPC server
// to group common settings and tasks inside single entity.
package grpcserver

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Server wraps gRPC server entity and handy means to simplify work with it.
type Server struct {
	address string
	server  *grpc.Server
	notify  chan error
}

// New creates new instance of gRPC server operating over SSL.
func New(address, crtPath, keyPath string, opts ...grpc.ServerOption) (*Server, error) {
	creds, err := credentials.NewServerTLSFromFile(crtPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("grpcserver - New - credentials.NewServerTLSFromFile: %w", err)
	}

	srvOpts := []grpc.ServerOption{
		grpc.Creds(creds),
	}
	srvOpts = append(srvOpts, opts...)

	grpcServer := grpc.NewServer(srvOpts...)

	s := &Server{
		address: address,
		server:  grpcServer,
		notify:  make(chan error, 1),
	}

	return s, nil
}

// Instance grants access to the underlying gRPC server.
// Should be used to attach new API services.
func (s *Server) Instance() *grpc.Server {
	return s.server
}

// Start launches the gRPC server.
func (s *Server) Start() {
	go func() {
		listen, err := net.Listen("tcp", s.address)
		if err != nil {
			s.notify <- err

			return
		}

		s.notify <- s.server.Serve(listen)
		close(s.notify)
	}()
}

// Notify reports errors received during start and work of the server.
// Usually such errors are not recoverable.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown() {
	if s.server == nil {
		return
	}

	s.server.GracefulStop()
}
