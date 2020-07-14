package servgrpc

import (
	"context"
	"github.com/marchmiel/proto-playground/Config"
	"github.com/marchmiel/proto-playground/proto"
	"github.com/marchmiel/proto-playground/server/wrapper"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func NewServePortChan(PORT string, Post chan<- wrapper.ClientDataType, sendback chan<- error) {
	lis, err := net.Listen("tcp", Config.GRPC_PORT) //React port + 2
	if err != nil {
		sendback <- errors.Wrap(err, "Could not setup listener")
	}
	gcrpServer := grpc.NewServer()
	proto.RegisterReservationServiceServer(gcrpServer, &server{Post})
	reflection.Register(gcrpServer)
	if err := gcrpServer.Serve(lis); err != nil {
		sendback <- errors.Wrap(err, "Could not server listener")
	}
}

type server struct {
	post chan<- wrapper.ClientDataType
}

func (s *server) MakeReservation(ctx context.Context, req *proto.BookTrip) (*proto.TripBooked, error) {
	clientPkg := NewGrpcData(req)
	s.post <- clientPkg
	err := <-clientPkg.err
	if err != nil {
		return nil, err
	}
	return clientPkg.respProto, nil
}
