package main

import (
	"github.com/marchmiel/proto-playground/client/model"
	"github.com/marchmiel/proto-playground/proto"
)

func (p *protos) BookTrip(mod *model.BookTripRequest) (*model.TripBookedResponse, error) {
	return &model.TripBookedResponse{}, nil
}

func 

type protos struct {
	slot proto.TripBooked
}

func main() {
	tos := protos{proto.TripBooked{}}

}
