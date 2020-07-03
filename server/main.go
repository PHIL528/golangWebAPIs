package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"proto-playground/Config"
	"proto-playground/proto"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var publisher_main message.Publisher

//publisher message.Publisher

func createTrip(client_name string) (*proto.TripBooked, error) { //This method is called by the PubSub puller and the gRPC handler
	booked_trip := proto.TripBooked{
		Trip: &proto.Trip{
			PassengerName: client_name,
			DriverName:    "Marek",
		},
	}
	jsonbytes, err := json.Marshal(booked_trip)
	if err != nil {
		//log.Printf("Could not convert json %v", err)
		return nil, errors.New("Could not convert JSON")
	}
	msg := message.NewMessage(watermill.NewUUID(), jsonbytes)
	if err := publisher_main.Publish(Config.Server_Publish_Topic, msg); err != nil {
		return nil, errors.New("Could not push confirmation to Pub/Sub service")
	}
	return &booked_trip, nil
}

func main() {
	os.Setenv("PUBSUB_EMULATOR_HOST", Config.Localhost_PubSub_PORT)

	//Generate Watermill Publisher
	logger := watermill.NewStdLogger(false, false)
	var err error
	publisher_main, err = googlecloud.NewPublisher(googlecloud.PublisherConfig{
		ProjectID: Config.PubSub_Project_Name,
	}, logger)
	if err != nil {
		panic(err)
	}

	//Generate Watermill Subscriber
	logger2 := watermill.NewStdLogger(false, false)
	subscriber, err := googlecloud.NewSubscriber(
		googlecloud.SubscriberConfig{
			GenerateSubscriptionName: func(topic string) string {
				return "server-puller"
			},
			ProjectID: Config.PubSub_Project_Name,
		},
		logger2,
	)
	if err != nil {
		panic(err)
	}
	messages, err := subscriber.Subscribe(context.Background(), Config.Server_Pull_Topic)
	if err != nil {
		panic(err)
	}

	go gRPCListener()      //To recieve gRPC requests from client via MakeReservation contract
	pubSubPuller(messages) //To recieve pulls from client on MakeReservation topic
}

//gRPC HANDLER

func gRPCListener() {
	lis, err := net.Listen("tcp", Config.GRPC_PORT) //React port + 2
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	gcrpServer := grpc.NewServer()
	proto.RegisterReservationServiceServer(gcrpServer, &server{log.New(os.Stdout, "gRPC Handler", log.LstdFlags)})
	reflection.Register(gcrpServer)
	if err := gcrpServer.Serve(lis); err != nil {
		log.Fatalf("Error %v", err)
	}
}

type server struct {
	l *log.Logger
}

func (s *server) MakeReservation(ctx context.Context, req *proto.BookTrip) (*proto.TripBooked, error) {
	s.l.Printf("Recieveing gRPC request to book trip from " + req.PassengerName)
	return createTrip(req.PassengerName)
}

//PubSub Handler

func pubSubPuller(messages <-chan *message.Message) {
	for msg := range messages {
		var BookTrip proto.BookTrip
		err := json.Unmarshal(msg.Payload, &BookTrip)
		if err != nil {
			fmt.Printf("Failed to unmarshal JSON")
		} else {
			tb, er := createTrip(BookTrip.PassengerName)
			if er != nil {
				fmt.Printf("Failed create trip")
			} else {
				fmt.Printf("PubSub reservation has been made by " + tb.Trip.PassengerName)
			}
		}
		log.Printf("received message: %s, payload: %s", msg.UUID, string(msg.Payload))
		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
	}
}
