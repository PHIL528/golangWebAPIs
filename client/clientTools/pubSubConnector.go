package clientTools

import (
	"errors"

	"github.com/marchmiel/proto-playground/Config"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"
	"golang.org/x/net/context"
)

//go run github.com/maxbrunsfeld/counterfeiter/v6 ./client/clientTools PubSubConnector
//counterfeiter:generate ./.. PubSubConnector
type PubSubConnector interface {
	SendReservation(*message.Message) (<-chan *message.Message, error)
}

func NewPubSubConnector() *pubSubConnector {
	return &pubSubConnector{}
}

type pubSubConnector struct {
	subscriptionName string
}

func (p *pubSubConnector) SendReservation(msg *message.Message) (<-chan *message.Message, error) {
	logger := watermill.NewStdLogger(false, false)
	subscriber, err := googlecloud.NewSubscriber(
		googlecloud.SubscriberConfig{
			GenerateSubscriptionName: func(topic string) string {
				return p.subscriptionName
			},
			ProjectID: Config.PubSub_Project_Name,
		},
		logger,
	)
	if err != nil {
		return nil, errors.New("Could not create subscriber")
	}
	messages, err := subscriber.Subscribe(context.Background(), Config.Server_Publish_Topic)
	if err != nil {
		return nil, errors.New("Could not create subscriber message channel")
	}
	//PUBLISHER
	publisher, err := googlecloud.NewPublisher(googlecloud.PublisherConfig{
		ProjectID: Config.PubSub_Project_Name,
	}, logger)
	if err != nil {
		return nil, errors.New("Could not create publisher")
	}
	if err := publisher.Publish(Config.Server_Pull_Topic, msg); err != nil {
		return nil, errors.New("Could not push request to Pub/Sub service")
	}
	return messages, nil
}
