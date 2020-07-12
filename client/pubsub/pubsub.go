package pubsub

import (
	"context"
	"encoding/json"
	//"errors"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/marchmiel/proto-playground/Config"
	"github.com/marchmiel/proto-playground/client/model"
	"github.com/pkg/errors"
	"time"
)

type pubSubTripBooker struct {
	publisher  message.Publisher
	subscriber message.Subscriber
	pubTopic   string
	subTopic   string
	subName    string
}

func (p *pubSubTripBooker) BookTrip(mod *model.BookTripRequest) (*model.TripBookedResponse, error) {
	defer p.Close()
	p.subName = mod.PassengerName
	stop, _ := context.WithTimeout(context.Background(), 10*time.Second)
	messages, err := p.subscriber.Subscribe(stop, Config.Server_Publish_Topic)
	if err != nil {
		return nil, errors.New("Could not create subscriber message channel")
	}
	msg, err := p.CreateMessage(mod)
	err = p.publisher.Publish(p.pubTopic, msg)
	if err != nil {
		return nil, errors.Wrap(err, "Could not publish message")
	}
	var template model.TripBookedResponse
	for {
		select {
		case <-stop.Done():
			return nil, stop.Err()
		case m := <-messages:
			if m.UUID == msg.UUID {
				e := json.Unmarshal(m.Payload, &template)
				if e != nil {
					return nil, e
				}
				return &template, nil
			} else {
				m.Ack()
				fmt.Printf("Recieved message from other UUID, resume research for correct UUID")
			}
		}
	}
}
func (p *pubSubTripBooker) Close() (error, error) { //errors not implemented yet
	p.subscriber.Close()
	p.publisher.Close()
	return nil, nil
}
func (p *pubSubTripBooker) CreateMessage(mod *model.BookTripRequest) (*message.Message, error) {
	jsonbytes, err := json.Marshal(mod)
	if err != nil {
		return nil, errors.Wrap(err, "Could not convert json")
	}
	return message.NewMessage(watermill.NewUUID(), jsonbytes), nil
}

//pub *message.Publisher, sub *message.Subscriber, pubTopic string, subTopic string
var NewTripBooker = func() (model.TripBooker, error) {
	template := pubSubTripBooker{}
	logger := watermill.NewStdLogger(false, false)
	var err error
	if template.publisher, err = googlecloud.NewPublisher(googlecloud.PublisherConfig{
		ProjectID: Config.PubSub_Project_Name,
	}, logger); err != nil {
		return nil, errors.Wrap(err, "Could not generate default Publisher")
	}
	if template.subscriber, err = googlecloud.NewSubscriber(
		googlecloud.SubscriberConfig{
			GenerateSubscriptionName: func(topic string) string {
				return "myasdf"
			},
			ProjectID: Config.PubSub_Project_Name,
		},
		logger,
	); err != nil {
		return nil, errors.Wrap(err, "Could not generate default Subscriber")
	}
	template.pubTopic = Config.Server_Pull_Topic
	template.subTopic = Config.Server_Publish_Topic

	return &template, nil //&pubSubTripBooker{publisher: pub, subscriber: sub, pubtopic: topicName, subtopic: subName}
}

/*
func GetEmptyPubSubTripBooker() *pubSubTripBooker {
	return &pubSubTripBooker{}
} */
func CustomPubSubTripBooker(pub *message.Publisher, sub *message.Subscriber, spubTopic string, ssubTopic string, ssubName string) model.TripBooker {
	return &pubSubTripBooker{publisher: *pub, subscriber: *sub, pubTopic: spubTopic, subTopic: ssubTopic, subName: ssubName}
}
