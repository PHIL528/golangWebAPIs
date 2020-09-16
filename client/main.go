package main

import (
	"fmt"
	"github.com/marchmiel/proto-playground/Config"
	"github.com/marchmiel/proto-playground/client/grpc"
	"github.com/marchmiel/proto-playground/client/model"
	"github.com/marchmiel/proto-playground/client/pubsub"
	"github.com/marchmiel/proto-playground/client/rest"
	"github.com/pkg/errors"
	"os"
	//"strings"
)

var pBTR *model.BookTripRequest
var pTBResp *model.TripBookedResponse

func main() {

	fmt.Println("Starting client")
	os.Setenv("PUBSUB_EMULATOR_HOST", Config.Localhost_PubSub_PORT)
	var tripBooker model.TripBooker
	var err error
	route := os.Args[1]
	clientName := os.Args[2]
	if route == "grpc" {
		tripBooker, err = grpc.NewTripBooker()
	} else if route == "pubsub" {
		tripBooker, err = pubsub.NewTripBooker()
	} else if route == "rest" {
		tripBooker, err = rest.NewTripBooker()
	} else {
		panic(errors.New("No valid route selected"))
	}
	if err != nil {
		panic(err)
	}
	bookTripRequest := model.NewBookTripRequest(clientName)
	tripBookedResponse, err := tripBooker.BookTrip(bookTripRequest)
	if err != nil {
		panic(errors.Wrap(err, "Could not create trip"))
	}
	fmt.Println("Assigned to driver " + tripBookedResponse.DriverName)

	pBTR = bookTripRequest
	pTBResp = tripBookedResponse
}
