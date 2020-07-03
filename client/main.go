package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/proto-playground/Config"
	"github.com/proto-playground/proto"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("client")
	route := strings.ToLower(os.Args[1])
	var err error
	var trip *proto.TripBooked
	if route == "grpc" {
		client_fname := strings.ToLower(os.Args[2])
		trip, err = send_via_gRPC(client_fname)
	} else if route == "pubsub" {
		client_fname := strings.ToLower(os.Args[2])
		trip, err = send_via_PubSub(client_fname)
	} else if route == "make" {
		time.Sleep(time.Second)
		fmt.Println("CLIENT: PLEASE ENTER IN ROUTE AND NAME INFORMATION")
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("CLIENT: Enter either 'grpc' or 'pubsub' : ")
		text1, _ := reader.ReadString('\n')
		fmt.Print("CLIENT: Enter first name of passenger : ")
		text2, _ := reader.ReadString('\n')
		if text1[:4] == "grpc" {
			fmt.Println("Sending by gRPC")
			trip, err = send_via_gRPC(text2)
		} else if text1[:6] == "pubsub" {
			fmt.Println("Sending bny PubSub")
			trip, err = send_via_PubSub(text2)
		} else {
			panic("Args 1 is neither gRPC nor PubSub")
		}
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

	logger := watermill.NewStdLogger(false, false)
	publisher, err := googlecloud.NewPublisher(googlecloud.PublisherConfig{
		ProjectID: Config.PubSub_Project_Name,
	}, logger)
	if err != nil {
		panic(err)
	}
	request := proto.BookTrip{
		PassengerName: client_name,
	}
	jsonbytes, err := json.Marshal(request)
	if err != nil {
		return nil, errors.New("Could not convert json: " + err.Error())
	}

	// PREPARED TO LISTEN FOR RESPONSE
	subscriber, err := googlecloud.NewSubscriber(
		googlecloud.SubscriberConfig{
			GenerateSubscriptionName: func(topic string) string {
				return client_name
			},
			ProjectID: Config.PubSub_Project_Name,
		},
		logger,
	)
	if err != nil {
		panic(err)
	}
	messages, err := subscriber.Subscribe(context.Background(), Config.Server_Publish_Topic)
	if err != nil {
		panic(err)
	}

	// PUBLISH REQUEST
	msg := message.NewMessage(watermill.NewUUID(), jsonbytes)
	if err := publisher.Publish(Config.Server_Pull_Topic, msg); err != nil {
		return nil, errors.New("Could not push request to Pub/Sub service")
	}

	// LISTEN FOR RESPONSE
	var TripBooked proto.TripBooked
	for msg := range messages {
		err := json.Unmarshal(msg.Payload, &TripBooked)
		if err != nil {
			fmt.Printf("Failed to unmarshal JSON")
		} else if TripBooked.Trip.PassengerName == client_name {
			log.Printf("Recieved trip confirmation, driver is " + TripBooked.Trip.DriverName)
			msg.Ack()
			break
		} else {
			log.Printf("Trip was from different client")
		}
		log.Printf("received message: %s, payload: %s", msg.UUID, string(msg.Payload))
		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
	}
	return &TripBooked, nil
}
