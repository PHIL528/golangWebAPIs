package grpc

import (
	"context"
	//"errors"
	"github.com/marchmiel/proto-playground/Config"
	"github.com/marchmiel/proto-playground/client/model"
	"github.com/marchmiel/proto-playground/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type grpcTripBooker struct {
	reservationClient proto.ReservationServiceClient
}

//res *proto.ReservationServiceClient
func NewTripBooker() (model.TripBooker, error) {
	con, err := grpc.Dial("localhost"+Config.GRPC_PORT, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "Could not create client connection")
	}
	res := proto.NewReservationServiceClient(con)
	return &grpcTripBooker{reservationClient: res}, nil
}

func (g *grpcTripBooker) BookTrip(mod *model.BookTripRequest) (*model.TripBookedResponse, error) {
	bookTripRequest := CreateProto(mod)
	ctx := context.Background()
	tripBooked, err := g.reservationClient.MakeReservation(ctx, bookTripRequest)
	return UnProto(tripBooked), errors.Wrap(err, "Failed to make reservation")
}
func CreateProto(mod *model.BookTripRequest) *proto.BookTrip {
	return &proto.BookTrip{
		PassengerName: mod.PassengerName,
	}
}
func UnProto(pro *proto.TripBooked) *model.TripBookedResponse {
	return &model.TripBookedResponse{
		PassengerName: pro.Trip.PassengerName,
		DriverName:    pro.Trip.DriverName,
	}
}

//New Trip boooker to setup booker, inject mock res client, in test invoke book trip with random data,
//should not return an error
//if empty return error
