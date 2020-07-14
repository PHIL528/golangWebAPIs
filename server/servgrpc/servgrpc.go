package servgrpc

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/marchmiel/proto-playground/client/model"
	"github.com/marchmiel/proto-playground/proto"
	//"github.com/marchmiel/proto-playground/server/wrapper"
)

type grpcDataType struct {
	reqProto  *proto.BookTrip
	respProto *proto.TripBooked
	err       chan error
}

func NewGrpcData(req *proto.BookTrip) *grpcDataType {
	return &grpcDataType{reqProto: req, err: make(chan error)}
}
func (g *grpcDataType) Unload(btr *model.BookTripRequest) error {
	btr.PassengerName = g.reqProto.PassengerName
	return nil
}
func (g *grpcDataType) Load(tbr *model.TripBookedResponse) error {
	g.respProto = &proto.TripBooked{
		Trip: &proto.Trip{
			PassengerName: tbr.PassengerName,
			DriverName:    tbr.DriverName,
		},
	}
	g.SendBack(nil)
	return nil
}
func (g *grpcDataType) CorrelationID() string {
	return watermill.NewUUID()
}
func (g *grpcDataType) SendBack(err error) {
	g.err <- err
}

//func (g *grpcDataType) GetResponse() interface{} {
//	return g.respProto
//}
