package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"

	"log"
	"os"
	"proto-playground/Config"
	"proto-playground/proto"
)

func main() {
	os.Setenv("PUBSUB_EMULATOR_HOST", Config.Localhost_PubSub_PORT)
	sub, _, err := Config.GetSubscriptionToTopic(context.Background(), Config.Server_Publish_Topic, "listener-pull", false)
	if err != nil {
		panic(err)
	}
	// MAINTAIN PULLING EVERY SECOND
	pull(sub)
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
		log.Fatalf("Error pull handler failed %v", err)
	}
}
