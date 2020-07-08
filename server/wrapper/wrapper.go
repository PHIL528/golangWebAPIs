package wrapper

import (
//"github.com/marchmiel/proto-playground/client/model"
)

type ClientDataTyp interface {
	Unload(*BookTripRequest) error
	Load(*TripBookedResponse) error
	Ret() interface{}
}
type BookTripRequest struct {
	PassengerName string
}
type TripBookedResponse struct {
	PassengerName string
	DriverName    string
}
