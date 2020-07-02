package main

import (
	"cloud.google.com/go/pubsub"
	"encoding/json"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"proto-playground/proto"
	"strings"
)

func main() {
	route := os.Args[1]
	client_fname := strings.ToLower(os.Args[2])
	if route == "grpc" {
		gRPC(client_fname)
	} else if route == "pubsub" {
		pubSub(client_fname)
	} else {
		log.Fatalf(route, " is neither gRPC nor PubSub")
	}
}
func gRPC(client_name string) {
	var con *grpc.ClientConn
	con, err := grpc.Dial("localhost:3002", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	defer con.Close()
	client := proto.NewReservationServiceClient(con)
	book_trip_request := proto.BookTrip{
		PassengerName: client_name,
	}
	confirmed_trip, err := client.MakeReservation(context.Background(), &book_trip_request)
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	log.Printf("Server assigned driver %s", confirmed_trip.DriverName)
}
func pubSub(client_name string) {
	os.Setenv("PUBSUB_EMULATOR_HOST", "localhost:8085")
	ps_ctx := context.Background()
	client, err := pubsub.NewClient(ps_ctx, "karhoo-local")
	if err != nil {
		log.Fatalf("Error making new client %v", err)
	}
	topic := client.Topic("events.MakeReservation")
	exist, err := topic.Exists(ps_ctx)
	if err != nil {
		log.Fatalf("Error in checking if topic exists %v", err)
	} else if !exist {
		log.Fatalf("events.MakeReservation topic has not been created")
	}
	defer topic.Stop()
	request := proto.BookTrip{
		PassengerName: client_name,
	}
	jsonbytes, err := json.Marshal(request)
	if err != nil {
		log.Printf("Could not convert json %v", err)
		return
	}
	ctx := context.Background() //Not sure if I should use a new context here.
	if _, err = topic.Publish(ctx, &pubsub.Message{Data: jsonbytes}).Get(ctx); err != nil {
		return
	}
}
