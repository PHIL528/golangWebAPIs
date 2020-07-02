package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"os"
	"proto-playground/Config"
	"proto-playground/proto"
	"strings"

	"cloud.google.com/go/pubsub"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	route := os.Args[1]
	client_fname := strings.ToLower(os.Args[2])
	var err error
	var trip *proto.TripBooked
	if route == "grpc" {
		trip, err = send_via_gRPC(client_fname)
	} else if route == "pubsub" {
		trip, err = send_via_PubSub(client_fname)
	} else {
		panic("Args 1 is neither gRPC nor PubSub")
	}
	if err == nil {
		fmt.Println("Assigned to driver " + trip.Trip.DriverName)
	} else {
		panic(err)
	}
}
func send_via_gRPC(client_name string) (*proto.TripBooked, error) {
	var con *grpc.ClientConn
	con, err := grpc.Dial("localhost:3002", grpc.WithInsecure())
	if err != nil {
		return &proto.TripBooked{}, err
	}
	defer con.Close()
	client := proto.NewReservationServiceClient(con)
	book_trip_request := proto.BookTrip{
		PassengerName: client_name,
	}
	confirmed_trip, err := client.MakeReservation(context.Background(), &book_trip_request)
	if err != nil {
		return &proto.TripBooked{}, err
	}
	return confirmed_trip, nil
	//log.Printf("Server assigned driver %s", confirmed_trip.DriverName)
}
func send_via_PubSub(client_name string) (*proto.TripBooked, error) {
	os.Setenv("PUBSUB_EMULATOR_HOST", Config.Localhost_PubSub_PORT)
	ps_ctx := context.Background()
	client, err := pubsub.NewClient(ps_ctx, Config.PubSub_Project_Name)
	if err != nil {
		return nil, errors.New("Error making new client " + err.Error())
	}
	send_topic := client.Topic(Config.Server_Pull_Topic)
	exist, err := send_topic.Exists(ps_ctx)
	if err != nil {
		return nil, errors.New("Error in checking if topic exists " + err.Error())
	} else if !exist {
		return nil, errors.New("events.MakeReservation topic has not been created")
	}
	defer send_topic.Stop()
	request := proto.BookTrip{
		PassengerName: client_name,
	}
	jsonbytes, err := json.Marshal(request)
	if err != nil {
		return nil, errors.New("Could not convert json: " + err.Error())
	}

	// PREPARED TO LISTEN FOR RESPONSE
	r_ctx := context.Background()
	sub, _, err := Config.GetSubscriptionToTopic(r_ctx, Config.Server_Publish_Topic, client_name, false)
	if err != nil {
		panic("Unable to recieve subscription to listen for confirmation")
	}

	// PUBLISH REQUEST
	if _, err = send_topic.Publish(ps_ctx, &pubsub.Message{Data: jsonbytes}).Get(ps_ctx); err != nil {
		return nil, errors.New("Publishing errors " + err.Error())
	}

	// LISTEN FOR RESPONSE
	b_ctx := context.Background()
	t_ctx, cancel := context.WithCancel(b_ctx)
	var TripBooked proto.TripBooked
	errx := sub.Receive(t_ctx, func(ctxxx context.Context, msg *pubsub.Message) {
		er := json.Unmarshal([]byte(msg.Data), &TripBooked)
		if er != nil {
			fmt.Printf("Could not unmarshal JSON, cannot confirm trip %v", er)
		} else {
			msg.Ack()
			cancel()
		}
	})
	if errx != nil {
		fmt.Println("Recieve failed")
	}
	return &TripBooked, nil
}
