package clientTools

import (
	"github.com/marchmiel/proto-playground/proto"
	"google.golang.org/grpc"
)

type ClientMaker interface {
	MakeClient() (proto.ReservationServiceClient, *grpc.ClientConn, error)
}

/*
type clientMaker struct {
	Conn        *grpc.ClientConn
	ResClient   proto.ReservationServiceClient
	HandlerType string
}
*/

//counterfeiter:generate ./../.. ClientMaker
