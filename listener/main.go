package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"
	"log"
	"os"
	"proto-playground/Config"
	"proto-playground/proto"
)

func main() {
	os.Setenv("PUBSUB_EMULATOR_HOST", Config.Localhost_PubSub_PORT)

	logs := watermill.NewStdLogger(false, false)
	subscriber, err := googlecloud.NewSubscriber(
		googlecloud.SubscriberConfig{
			GenerateSubscriptionName: func(topic string) string {
				return "listener"
			},
			ProjectID: Config.PubSub_Project_Name,
		},
		logs,
	)
	if err != nil {
		panic(err)
	}
	messages, err := subscriber.Subscribe(context.Background(), Config.Server_Publish_Topic)
	if err != nil {
		panic(err)
	}
	// MAINTAIN PULLING EVERY SECOND
	pull(messages)
}

func pull(messages <-chan *message.Message) {
	for msg := range messages {
		var TripBooked proto.TripBooked
		err := json.Unmarshal(msg.Payload, &TripBooked)
		if err != nil {
			fmt.Printf("Failed to unmarshal JSON")
		} else {
			log.Printf("Logging trip made by %v", TripBooked.Trip.PassengerName)
			log.Printf("To be serviced by %v", TripBooked.Trip.DriverName)
		}
		log.Printf("received message: %s, payload: %s", msg.UUID, string(msg.Payload))
		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
	}
}
