package wrapper

import (
	"github.com/marchmiel/proto-playground/client/model"
)

type ClientDataType interface {
	Unload(*model.BookTripRequest) error
	Load(*model.TripBookedResponse) error
	CorrelationID() string
	SendBack(error)
}
