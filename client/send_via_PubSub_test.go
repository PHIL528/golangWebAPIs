package main

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	//"github.com/marchmiel/proto-playground/clientTools"
	"fmt"
	"github.com/marchmiel/proto-playground/conv"
	"github.com/marchmiel/proto-playground/proto"
	mock "github.com/marchmiel/proto-playground/protoplaygroundfakes"
	"testing"
)

var testPassenger string = "PhilMur"

func MapPubSubInputToOutput(msg *message.Message) (<-chan *message.Message, error) {
	var BookTrip proto.BookTrip
	json.Unmarshal(msg.Payload, &BookTrip)
	TripBooked := proto.TripBooked{
		Trip: &proto.Trip{
			PassengerName: BookTrip.PassengerName,
			DriverName:    "TestDriverMarek",
		},
	}
	Msg, _ := conv.ToJsonToMessage(TripBooked)
	pubSubOut := make(chan *message.Message)
	go func(sendTo chan<- *message.Message) {
		sendTo <- Msg
	}(pubSubOut)
	return pubSubOut, nil
}

func TestSendViaPubSub(t *testing.T) {
	fmt.Println("TESTING")
	VarPubSubConnector := mock.FakePubSubConnector{}
	VarPubSubConnector.SetConnType("FakeConnector")
	VarPubSubConnector.SendReservationStub = MapPubSubInputToOutput

	BookedTrip, _ := send_via_PubSub(testPassenger)
	fmt.Println("BOOKED TRIP")
	fmt.Println(BookedTrip)
}

//go test client/send_via_gRPC_test.go -v
