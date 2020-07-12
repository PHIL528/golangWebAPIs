package servpubs

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/marchmiel/proto-playground/client/mocks"
	"github.com/marchmiel/proto-playground/client/model"
	"github.com/marchmiel/proto-playground/server/wrapper"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var testPassenger string = "testPassenger505"

func TestWithMockSubscriber(t *testing.T) {
	var testBookRequest *model.BookTripRequest = &model.BookTripRequest{
		PassengerName: testPassenger,
	}
	jsonbytes, _ := json.Marshal(testBookRequest)
	testMsg := message.NewMessage(watermill.NewUUID(), jsonbytes)
	testChan := make(chan *message.Message, 20)
	testChan <- testMsg

	generateSubscriber = func() (message.Subscriber, error) {
		subscriber := &mocks.FakeSubscriber{}
		subscriber.SubscribeReturns(testChan, nil)
		return subscriber, nil
	}

	collection := make(chan wrapper.ClientDataType, 100)
	errorsChan := make(chan error, 3)
	go NewServePortChan("nil", collection, errorsChan)

	var invoke wrapper.ClientDataType
	select {
	case <-time.After(time.Second * 5):
		t.Error("Timeout on recieving from collection")
	case invk := <-collection:
		invoke = invk
	}
	var btrInvoked model.BookTripRequest
	invoke.Unload(&btrInvoked)
	assert.Equal(t, testBookRequest, &btrInvoked, " comparing sent request with unloaded request")

}
