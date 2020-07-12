package servgrpc

import (
	"context"
	"github.com/marchmiel/proto-playground/client/model"
	"github.com/marchmiel/proto-playground/proto"
	"github.com/marchmiel/proto-playground/server/wrapper"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var testPassenger string = "testPassenger606"

func TestGrpc(t *testing.T) {
	testRequest := &proto.BookTrip{
		PassengerName: "testPassenger606",
	}
	expectedTripBookRequest := &model.BookTripRequest{
		PassengerName: "testPassenger606",
	}
	collection := make(chan wrapper.ClientDataType, 50)
	gServ := server{collection}
	gReturn := make(chan *proto.TripBooked)
	go func() {
		tripBooked, _ := gServ.MakeReservation(context.Background(), testRequest)
		gReturn <- tripBooked
	}()

	var invoke wrapper.ClientDataType
	select {
	case <-time.After(time.Second * 5):
		t.Error("Timeout on recieving from collection")
	case invk := <-collection:
		invoke = invk
	}
	var bookTripRequest model.BookTripRequest
	invoke.Unload(&bookTripRequest)
	assert.Equal(t, expectedTripBookRequest, &bookTripRequest, " comparing collected tripRequest with expected")

	StubResponse := model.TripBookedResponse{
		PassengerName: "testPassenger606",
		DriverName:    "TestDriver606",
	}
	expectedResponse := &proto.TripBooked{
		Trip: &proto.Trip{
			PassengerName: "testPassenger606",
			DriverName:    "TestDriver606",
		},
	}
	invoke.Load(&StubResponse)
	select {
	case <-time.After(time.Second * 5):
		t.Error("Timeout recieved on returning gRPC")
	case resp := <-gReturn:
		assert.Equal(t, expectedResponse, resp, " comparing server's return of StubResponse to expectedResponse")
	}
}
