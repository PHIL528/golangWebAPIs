package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"proto-playground/Config"
	"proto-playground/proto"
	"time"
)

func main() {
	os.Setenv("PUBSUB_EMULATOR_HOST", Config.Localhost_PubSub_PORT)

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, Config.PubSub_Project_Name)
	if err != nil {
		fmt.Println("Could not setup new client")
	}
	topic := client.Topic(Config.Server_Publish_Topic) //it is the server that publishes, not this
	if exists, err := topic.Exists(ctx); !exists {
		log.Fatalf("Topic does not exist %v", err)
	}
	if err != nil {
		log.Fatal("Other error with topic checking")
	}
	sub := client.Subscription("listener_pull")
	exists, err := sub.Exists(ctx)
	if err != nil {
		log.Fatal("Error checking if sub exists")
	}
	if !exists {
		sub, err = client.CreateSubscription(ctx, "listener_pull", pubsub.SubscriptionConfig{Topic: topic})
		if err != nil {
			log.Fatalf("Cannot create/join subscription")
		}
	}
	// MAINTAIN PULLING EVERY SECOND
	for {
		time.Sleep(time.Second)
		pull(sub)
	}
}
func pull(s *pubsub.Subscription) {
	ctx := context.Background()
	//	client, _ := pubsub.NewClient(ctx, "karhoo-local")
	err := s.Receive(ctx, func(ctxx context.Context, msg *pubsub.Message) {
		var TripBooked proto.TripBooked
		er := json.Unmarshal([]byte(msg.Data), &TripBooked)
		if er != nil {
			log.Printf("Could not unmarshal JSON, cannot confirm trip %v", er)
		} else {
			log.Printf("Logging trip made by %v", TripBooked.Trip.PassengerName)
			log.Printf("To be serviced by %v", TripBooked.Trip.DriverName)
			msg.Ack()
		}
	})
	if err != nil {
		log.Fatalf("Error  %v", err)
	}
}
