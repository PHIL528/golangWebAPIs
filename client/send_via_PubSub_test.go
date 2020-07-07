package main

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	//"github.com/marchmiel/proto-playground/clientTools"
	"fmt"
	"github.com/marchmiel/proto-playground/conv"
	"github.com/marchmiel/proto-playground/proto"
	mock "github.com/marchmiel/proto-playground/protoplaygroundfakes"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testPassenger string = "PhilMur"
var ExpectedTrip *proto.TripBooked = &proto.TripBooked{
	Trip: &proto.Trip{
		PassengerName: testPassenger,
		DriverName:    "TestDriverMarek",
	},
}

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
	t.Logf("Starting PubSub Test")
	VarPubSubConnector := mock.FakePubSubConnector{}
	VarPubSubConnector.SetConnType("FakeConnector")
	VarPubSubConnector.SendReservationStub = MapPubSubInputToOutput

	BookedTrip, err := send_via_PubSub(testPassenger)
	fmt.Println("is nil?")
	fmt.Println(BookedTrip)
	fmt.Println(err)
	assert.Equal(t, ExpectedTrip, BookedTrip, "Comparing expected and return trip")
}

//go test client/send_via_gRPC_test.go -v
