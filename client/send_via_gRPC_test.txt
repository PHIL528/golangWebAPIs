package main

import (
	//"github.com/marchmiel/proto-playground/clientTools"
	"context"
	"fmt"
	"github.com/marchmiel/proto-playground/proto"
	"github.com/marchmiel/proto-playground/proto/protofakes"
	mock "github.com/marchmiel/proto-playground/protoplaygroundfakes"
	"google.golang.org/grpc"
	"testing"
)

var testPassengerX string = "PhilMur"

func MapGRPCInputToOutput(ctx context.Context, BookTrip *proto.BookTrip, opts ...grpc.CallOption) (*proto.TripBooked, error) {
	return &proto.TripBooked{
		Trip: &proto.Trip{
			PassengerName: BookTrip.PassengerName,
			DriverName:    "Test driver Marek",
		},
	}, nil
}
func TestViaGRPC(t *testing.T) {

	//fCon := mock.FakeClientConnInterface{}
	fRes := &protofakes.FakeReservationServiceClient{}
	fRes.MakeReservationStub = MapGRPCInputToOutput
	VarClientMaker := mock.FakeClientMaker{}
	VarClientMaker.MakeClientReturns(fRes, nil, nil)

	BookedTrip, _ := send_via_PubSub(testPassengerX)
	fmt.Println("BOOKED TRIP")
	fmt.Println(BookedTrip)
}

//go test -v client/send_via_gRPC_test.go
