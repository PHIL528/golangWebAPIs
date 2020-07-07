package wrapper

import (
	"github.com/marchmiel/proto-playground/client/model"
)

type ClientDataType interface {
	Unload(*model.BookTripRequest) (*model.BookTripRequest, error)
	Encode(*model.TripBookedResponse) error
}
