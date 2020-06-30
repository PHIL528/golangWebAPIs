package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"log"
	//"math/rand"
	"os"
	"proto-playground/proto"
	"time"
)

func main() {
	fmt.Println("Starting listener")
	os.Setenv("PUBSUB_EMULATOR_HOST", "localhost:8085")

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, "karhoo-local")
	topic := client.Topic("events.TripBooked")
	if exists, err := topic.Exists(ctx); !exists {
		log.Fatalf("Topic does not exist %v", err)
	}
	sub, err := client.CreateSubscription(ctx, "yeet", pubsub.SubscriptionConfig{Topic: topic}) //I tried adding an Endpoint URL in the PushConfig but I don't think it works without a certificate
	if err != nil {
		fmt.Println("Subscription may already exist, attempting to join")
		sub = client.Subscription("yeet")
		if sub == nil {
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
