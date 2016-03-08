package receive

import (
	//	"errors"
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	api "github.com/njpatel/portal/api"
	"github.com/njpatel/portal/shared"
)

// Run starts a receiving client
func Run(c *Config, token string, outputDir string) {
	client := connect(c)

	// Authentication happens by sending the secret token via metadata
	md := metadata.Pairs(shared.SecretKey, c.Secret)

	stream, err := client.Get(metadata.NewContext(context.Background(), md), &api.GetRequest{
		Token: token,
	})
	shared.ExitOnError(err, "Unable to initiate Receive: %v", grpc.ErrorDesc(err))

	for {
		res, err := stream.Recv()
		if err != nil {
			shared.ExitOnError(err, "Unable to receive data: %v", grpc.ErrorDesc(err))
		}
		fmt.Println("received", res.Type)
	}
}

func connect(c *Config) api.PortalClient {
	var gOpts []grpc.DialOption
	if c.Insecure == true {
		gOpts = append(gOpts, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(c.Address, gOpts...)
	shared.ExitOnError(err, "Unable to connect to Portal server %v: %v", c.Address, grpc.ErrorDesc(err))

	return api.NewPortalClient(conn)
}
