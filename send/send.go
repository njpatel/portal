package send

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	api "github.com/njpatel/portal/api"
	"github.com/njpatel/portal/shared"
)

// Config describes the configuration of the client
type Config struct {
	Address  string
	Insecure bool
	Secret   string
}

// Run starts a new client for one-shot sending
func Run(c *Config, args []string) {
	run(c, false, args)
}

// RunSync starts a new client for sync sending
func RunSync(c *Config, args []string) {
	run(c, true, args)
}

func run(c *Config, sync bool, args []string) {
	client := connect(c)

	// Authentication happens by sending the secret token via metadata
	md := metadata.Pairs(shared.SecretKey, c.Secret)

	stream, err := client.Put(metadata.NewContext(context.Background(), md))
	shared.ExitOnError(err, "Unable to initiate Put: %v", grpc.ErrorDesc(err))

	// We should be told to wait & be given a token for the receive stream
	res, err := stream.Recv()
	shared.ExitOnError(err, "Unable to begin session: %v", grpc.ErrorDesc(err))

	if res.Type != api.SessionStateType_WAIT || res.Token == "" {
		err = errors.New("Incorrect session response received")
		shared.ExitOnError(err, err.Error())
	}

	fmt.Printf("Session sucessful. Receive token is %s\n", res.Token)

	// This blocks until a receiver connects with the same token, or a timeout is reached
	res, err = stream.Recv()
	shared.ExitOnError(err, "Unable to commence sending data: %v", grpc.ErrorDesc(err))

	if res.Type != api.SessionStateType_ACTIVE || res.ReceiverIP == "" {
		err = errors.New("Incorrect session commence response received")
		shared.ExitOnError(err, err.Error())
	}

	fmt.Printf("Session connected: ReceiverIP = %s. Commencing transfers\n", res.ReceiverIP)

	go receive(stream)

	ticker := time.NewTicker(time.Second * 5)
	for now := range ticker.C {
		fmt.Println("ping", now)
		err := stream.Send(&api.Frame{
			Type: api.FrameType_PING,
		})
		shared.ExitOnError(err, "Canceled: %v", err)
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

func receive(stream api.Portal_PutClient) {
	for {
		res, err := stream.Recv()
		if err == nil && res.Type == api.SessionStateType_CANCEL {
			err = errors.New(res.Reason)
		}
		shared.ExitOnError(err, grpc.ErrorDesc(err))
	}
}
