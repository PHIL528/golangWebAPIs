package model

import ()

type TripBooker interface {
	BookTrip(BookTripRequest) error
}
type BookTripRequest struct {
	PassengerName string
	Destination   string
}
