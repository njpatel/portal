package server

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	api "github.com/njpatel/portal/api"
)

// Service ...
type service struct {
	config *Config
}

// Run ...
func Run(c *Config) {
	s := &service{
		config: c,
	}

	listener, err := net.Listen("tcp", c.Address)
	if err != nil {
		fmt.Printf("Error binding to %v: %s\n", c.Address, err)
		return
	}
	defer listener.Close()

	server := grpc.NewServer()
	api.RegisterPortalServer(server, s)
	fmt.Printf("Started portal server on %v:%v\n", listener.Addr().Network(), listener.Addr().String())
	server.Serve(listener)
}

func (*service) Put(stream api.Portal_PutServer) error {
	return nil
}

func (*service) Get(req *api.GetRequest, stream api.Portal_GetServer) error {
	return nil
}
