package clientTools

import (
	"github.com/marchmiel/proto-playground/proto"
	"google.golang.org/grpc"
)

type ClientMaker interface {
	MakeClient() (proto.ReservationServiceClient, *grpc.ClientConn, error)
}

//go run github.com/maxbrunsfeld/counterfeiter/v6 ./clientTools/grpcMock.go ClientMaker

/*
type clientMaker struct {
	Conn        *grpc.ClientConn
	ResClient   proto.ReservationServiceClient
	HandlerType string
}
*/

//counterfeiter:generate ./../.. ClientMaker
