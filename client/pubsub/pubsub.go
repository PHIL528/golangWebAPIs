package pubsub

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/marchmiel/proto-playground/client/model"
)

type pubSubTripBooker struct {
	publisher message.Publisher
	topic     string
}

func NewTripBooker(pub message.Publisher, topicName string) model.TripBooker {
	return &pubSubTripBooker{publisher: pub, topic: topicName}
}
func (p *pubSubTripBooker) BookTrip(mod model.BookTripRequest) error {
	msg := p.CreateMessage(mod)
	err := p.publisher.Publisher(p.topic, msg)
	if err != nil {
		err = errors.Wrap(err, "Could not publish message")
	}
	return err
}
func (p *pubSubTripBooker) CreateMessage(mod model.BookTripRequest) (*message.Message, error) {
	jsonbytes, err := json.Marshal(mod)
	if err != nil {
		return nil, errors.Wrap(err, "Could not convert json")
	}
	return message.NewMessage(watermill.NewUUID(), jsonbytes), nil
}
