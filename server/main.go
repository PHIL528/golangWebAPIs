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
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var publish_topic *pubsub.Topic

func createTrip(client_name string) (*proto.TripBooked, error) { //This method is called by the PubSub puller and the gRPC handler
	booked_trip := proto.TripBooked{
		Trip: &proto.Trip{
			PassengerName: client_name,
			DriverName:    "Marek",
		},
	}
	jsonbytes, err := json.Marshal(booked_trip)
	if err != nil {
		log.Printf("Could not convert json %v", err)
		return &proto.TripBooked{}, errors.New("Could not convert JSON")
	}
	ps_ctx := context.Background()
	if _, err = publish_topic.Publish(ps_ctx, &pubsub.Message{Data: jsonbytes}).Get(ps_ctx); err != nil {
		return &proto.TripBooked{}, errors.New("Could not push confirmation to Pub/Sub service")
	}
	return &booked_trip, nil
}

func main() {
	publish_topic := setupPublisher()  //To publish confirmed requests, to be recieved by client + listener
	pull_subscription := setupPuller() //To recieve requests from client

	go gRPCListener(publish_topic)  //To recieve gRPC requests from client via MakeReservation contract
	pubSubPuller(pull_subscription) //To recieve pulls from client on MakeReservation topic
}

//gRPC HANDLER

func gRPCListener(publish_topic *pubsub.Topic) {
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

func pubSubPuller(sub *pubsub.Subscription) {
	puller_log := log.New(os.Stdout, "Server/PubSub-Puller: ", log.LstdFlags)
	ctx := context.Background()
	err := sub.Receive(ctx, func(ctxx context.Context, msg *pubsub.Message) {
		var BookTrip proto.BookTrip
		er := json.Unmarshal([]byte(msg.Data), &BookTrip)
		if er != nil {
			puller_log.Printf("Failed to unmarshal JSON")
		} else {
			tb, e := createTrip(BookTrip.PassengerName)
			puller_log.Printf("PubSub reservation has been made by " + tb.Trip.PassengerName)
			if e != nil {
				puller_log.Printf("Failed create trip")
			}
			msg.Ack()
		}
	})
	puller_log.Fatalf("Handler failed " + err.Error())
}

//SETUP

func setupPublisher() *pubsub.Topic {
	os.Setenv("PUBSUB_EMULATOR_HOST", Config.Localhost_PubSub_PORT)
	topic, _, err := Config.GetTopic(context.Background(), Config.Server_Publish_Topic, true)
	if err != nil {
		panic(err)
	}
	return topic
}

func setupPuller() *pubsub.Subscription {
	sub, _, err := Config.GetSubscriptionToTopic(context.Background(), Config.Server_Pull_Topic, "server-pull", true)
	if err != nil {
		panic(err)
	}
	return sub
}
