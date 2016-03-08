package server

import (
	"fmt"
	"net"
	"strings"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	api "github.com/njpatel/portal/api"
	"github.com/njpatel/portal/shared"
)

var errf = grpc.Errorf // Work around `go vet` madness

// Service ...
type service struct {
	config   *Config
	exchange *exchange
}

// Run ...
func Run(c *Config) {
	s := &service{
		config:   c,
		exchange: newExchange(),
	}

	listener, err := net.Listen("tcp", c.Address)
	if err != nil {
		fmt.Printf("Error binding to %v: %s\n", c.Address, err)
		return
	}
	defer listener.Close()

	server := grpc.NewServer()
	api.RegisterPortalServer(server, s)
	fmt.Printf("Started portal server on %v://%v\n", listener.Addr().Network(), listener.Addr().String())
	server.Serve(listener)
}

func (s *service) Put(stream api.Portal_PutServer) error {
	addr := remoteAddrFromContext(stream.Context())
	logger.Tracef("Received new Put request from %s", addr)

	if err := authenticate(stream.Context(), s.config.Secret); err != nil {
		return err
	}

	// Create the session, and send back the token
	token, waitChan, sendChan := s.exchange.newSession(addr)
	err := stream.Send(&api.PutResponse{
		Type:  api.SessionStateType_WAIT,
		Token: token,
	})
	if err != nil {
		logger.Warningf("Unable to send token: %v", err)
		return s.cleanup(token, errf(codes.Internal, "Canceled"))
	}

	err = s.waitForReceiver(token, waitChan, stream)
	if err != nil {
		return s.cleanup(token, err)
	}

	// Forward incoming frames to our receiver
	for {
		frame, err := stream.Recv()
		if err != nil {
			if strings.Contains(err.Error(), "context canceled") == false {
				logger.Warningf("Unable to receive Put: %v", err)
			}
			return s.cleanup(token, errf(codes.Aborted, "Canceled"))
		}

		select {
		// If waitChan closes, our receiver has gone missing
		case _, ok := <-waitChan:
			if ok == false {
				stream.Send(&api.PutResponse{
					Type:   api.SessionStateType_CANCEL,
					Reason: "Receiver disconnected",
				})
				return errf(codes.Aborted, "Canceled")
			}

		// Otherwise just foward
		default:
			if frame.Type != api.FrameType_PING {
				sendChan <- frame
			}
		}
	}
}

func (s *service) waitForReceiver(token string, waitChan <-chan string, stream api.Portal_PutServer) error {
	select {
	case <-time.After(1 * time.Minute):
		return errf(codes.DeadlineExceeded, "No receive request received for 60 seconds")

	case receiverIP, ok := <-waitChan:
		if receiverIP == "" || ok == false {
			return errf(codes.Aborted, "Unable to initiate receive")
		}
		// Let the sender know to start sending
		err := stream.Send(&api.PutResponse{
			Type:       api.SessionStateType_ACTIVE,
			ReceiverIP: receiverIP,
		})
		if err != nil {
			return errf(codes.Aborted, "Canceled")
		}
	}
	return nil
}

func (s *service) Get(req *api.GetRequest, stream api.Portal_GetServer) error {
	addr := remoteAddrFromContext(stream.Context())
	logger.Tracef("Received new Get request from %s", addr)

	if err := authenticate(stream.Context(), s.config.Secret); err != nil {
		return err
	}

	token := req.Token
	senderChan, err := s.exchange.connect(token, addr)

	for {
		select {
		case <-stream.Context().Done():
			return s.cleanup(token, errf(codes.Aborted, "Canceled"))

		case frame, ok := <-senderChan:
			if ok == false {
				return s.cleanup(token, errf(codes.Aborted, "Canceled"))
			}
			err = stream.Send(frame)
			if err != nil {
				if strings.Contains(err.Error(), "context canceled") == false {
					logger.Warningf("Unable to send data: %v", err)
				}
				return s.cleanup(token, errf(codes.Aborted, "Send error: %v", err))
			}
		}
	}
}

func (s *service) cleanup(token string, err error) error {
	s.exchange.delete(token)
	return err
}

func remoteAddrFromContext(ctx context.Context) string {
	peer, ok := peer.FromContext(ctx)
	if ok == true {
		return strings.Split(peer.Addr.String(), ":")[0]
	}
	return "Unknown"
}

func authenticate(ctx context.Context, serverSecret string) error {
	md, ok := metadata.FromContext(ctx)
	if ok != true {
		logger.Warningf("Unable to get connection metadata")
		return errf(codes.Unauthenticated, "Missing metadata in context")
	}

	secret := md[shared.SecretKey]
	if len(secret) != 1 || secret[0] != serverSecret {
		logger.Tracef("Rejecting secret. Client=%s Server=%s", secret, serverSecret)
		return errf(codes.Unauthenticated, "Unauthorized - Bad secret")
	}

	return nil
}
