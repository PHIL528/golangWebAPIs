package model

import ()

type TripBooker interface {
	BookTrip(*BookTripRequest) (*TripBookedResponse, error)
}
type BookTripRequest struct {
	PassengerName string
}
type TripBookedResponse struct {
	PassengerName string
	DriverName    string
}

func NewBookTripRequest(name string) *BookTripRequest {
	return &BookTripRequest{PassengerName: name}
}
func NewTripBooked(t *BookTripRequest, d string) *TripBookedResponse {
	return &TripBookedResponse{PassengerName: t.PassengerName, DriverName: d}
}
