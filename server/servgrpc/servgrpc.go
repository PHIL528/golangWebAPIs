package servgrpc

import (
	//"github.com/marchmiel/proto-playground/client/model"
	"github.com/marchmiel/proto-playground/proto"
	"github.com/marchmiel/proto-playground/server/wrapper"
)

type grpcData struct {
	ReqProto  *proto.BookTrip
	RespProto *proto.TripBooked
}

func NewGrpcData(req *proto.BookTrip) *grpcData {
	return &grpcData{ReqProto: req}
}

func (g *grpcData) Unload(btr *wrapper.BookTripRequest) error {
	btr.PassengerName = g.ReqProto.PassengerName
	return nil
}
func (g *grpcData) Load(tbr *wrapper.TripBookedResponse) error {
	g.RespProto = &proto.TripBooked{
		Trip: &proto.Trip{
			PassengerName: tbr.PassengerName,
			DriverName:    tbr.DriverName,
		},
	}
	return nil
}
func (g *grpcData) Ret() interface{} {
	return nil
}
