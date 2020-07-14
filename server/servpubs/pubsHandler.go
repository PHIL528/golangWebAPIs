package servpubs

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/marchmiel/proto-playground/Config"
	"github.com/marchmiel/proto-playground/server/wrapper"
	"github.com/pkg/errors"
	"log"
	"os"
)

type pubsHandler struct {
	post chan<- wrapper.ClientDataType
}

func NewServePortChan(PORT string, Post chan<- wrapper.ClientDataType, sendback chan<- error) {
	os.Setenv("PUBSUB_EMULATOR_HOST", Config.Localhost_PubSub_PORT)
	subscriber, err := generateSubscriber()
	if err != nil {
		sendback <- errors.Wrap(err, "Could not create subscriber")
	}
	messages, err := subscriber.Subscribe(context.Background(), Config.Server_Pull_Topic)
	if err != nil {
		sendback <- errors.Wrap(err, "Could not create reciever messages channel")
	}
	handler := pubsHandler{Post}
	for msg := range messages {
		go handler.pubsHandle(msg)
	}
}

var generateSubscriber = func() (message.Subscriber, error) {
	return googlecloud.NewSubscriber(
		googlecloud.SubscriberConfig{
			GenerateSubscriptionName: func(topic string) string {
				return "server-puller"
			},
			ProjectID: Config.PubSub_Project_Name,
		},
		watermill.NewStdLogger(false, false),
	)
}

func (p *pubsHandler) pubsHandle(msg *message.Message) {
	clientPkg := NewPubsData(msg)
	p.post <- clientPkg
	err := <-clientPkg.err
	if err != nil {
		log.Printf(errors.Wrap(err, " pubsHandler failed").Error())
		msg.Nack()
	} else {
		msg.Ack()
	}
}
