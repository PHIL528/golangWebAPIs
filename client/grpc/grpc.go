package grpc

import (
	"errors"
	"github.com/marchmiel/proto-playground/client/model"
	"github.com/marchmiel/proto-playground/proto"
)

type grpcTripBooker struct {
	reservationClient proto.ReservationServiceClient
}

func NewTripBooker(res proto.ReservationServiceClient) model.TripBooker {
	return &gripTripBooker{reservationClient: res}
}

func (g *grpcTripBooker) BookTrip(req model.BookTripRequest) error {
	bookTripRequest := proto.BookTrip{
		PassengerName: req.PassengerName,
	}
	tripBooked, err := g.reservationClient.MakeReservation(context.Background(), &bookTripRequest)
	fmt.Println(tripBooked)
	return errors.Wrap(err, "Failed to make reservation")
}

//New Trip boooker to setup booker, inject mock res client, in test invoke book trip with random data,
//should not return an error
//if empty return error

//million dollar mistake
