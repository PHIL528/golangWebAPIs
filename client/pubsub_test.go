package main

//go run github.com/maxbrunsfeld/counterfeiter/v6 ./client/model TripBooker
//go run github.com/maxbrunsfeld/counterfeiter/v6 github.com/ThreeDotsLabs/watermill/message.Publisher
//go run github.com/maxbrunsfeld/counterfeiter/v6 github.com/ThreeDotsLabs/watermill/message.Subscriber

import (
	"context"
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/marchmiel/proto-playground/Config"
	"github.com/marchmiel/proto-playground/client/mocks"
	"github.com/marchmiel/proto-playground/client/model"
	"github.com/marchmiel/proto-playground/client/pubsub"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	//"log"
)

var testPassenger string = "PhilMur"

//Noedit
var testArgs = []string{"cmd", "pubsub", testPassenger}
var ExpectedRequest *model.BookTripRequest = &model.BookTripRequest{
	PassengerName: testPassenger,
}
var ExpectedTrip *model.TripBookedResponse = &model.TripBookedResponse{
	PassengerName: testPassenger,
	DriverName:    "TestDriver101PubSub",
}

func TestWithMockPubSub(t *testing.T) {
	t.Log("Testing client with pubsub")
	pubsub.NewTripBooker = MockTripBooker
	os.Args = []string{"cmd", "pubsub", testPassenger}
	main()
	assert.Equal(t, pBTR, ExpectedRequest, " comparing trip requests")
	assert.Equal(t, pTBResp, ExpectedTrip, " comparing trip requests")
}

var MockTripBooker = func() (model.TripBooker, error) {
	MockChan = make(chan *message.Message, 10)
	publisher := func() message.Publisher {
		publisher := &mocks.FakePublisher{}
		publisher.PublishStub = PublishStub
		return publisher
	}()
	subscriber := func() message.Subscriber {
		subscriber := &mocks.FakeSubscriber{}
		subscriber.SubscribeStub = SubscribeStub
		return subscriber
	}()
	template := pubsub.CustomPubSubTripBooker(&publisher, &subscriber, Config.Server_Pull_Topic, Config.Server_Publish_Topic, "myasdf")
	return template, nil
}

var PublishStub = func(top string, ms ...*message.Message) error { //Acts as client sender and server publisher
	msg := ms[0]
	var BTR model.BookTripRequest
	json.Unmarshal(msg.Payload, &BTR)
	TBResp := model.TripBookedResponse{
		PassengerName: BTR.PassengerName,
		DriverName:    "TestDriver101PubSub",
	}
	jsonbytes, _ := json.Marshal(TBResp)
	MockChan <- message.NewMessage(msg.UUID, jsonbytes)
	return nil
}
var MockChan chan *message.Message
var SubscribeStub = func(ctx context.Context, topic string) (<-chan *message.Message, error) {
	return MockChan, nil
}
