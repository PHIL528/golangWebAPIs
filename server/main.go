package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"proto-playground/proto"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var topic *pubsub.Topic

func main() {
	fmt.Println("Starting server")

	//CONNECTING TO PUBSUB
	fmt.Println("Starting PubSub connection")
	os.Setenv("PUBSUB_EMULATOR_HOST", "localhost:8085")
	ps_ctx := context.Background()
	notified := false
	var client *pubsub.Client
	var err error
	for { //I put the for loop in hear to loop until the PubSub is started, but its useless because err returns nil regardless of whether or not the PubSub terminal is started
		client, err = pubsub.NewClient(ps_ctx, "karhoo-local")
		if err == nil {
			break
		} else if !notified {
			log.Printf("Failed to create pubsub client, %v", err)
			log.Printf("Perhaps the PubSub terminal has not yet been started, will reattempt conncetion once per second")
			notified = true
		}
		time.Sleep(time.Second)
	}
	topic, err = client.CreateTopic(ps_ctx, "events.TripBooked")
	if err != nil {
		log.Printf("Failed to create topic %v", err)
		log.Printf("Perhaps the topic already exists, joining topic instead of creating")
		topic = client.Topic("events.TripBooked")
		if exists, er := topic.Exists(ps_ctx); !exists {
			log.Fatalf("Cannot create topic and topic does not exist %v", er)
		}
	}
	defer topic.Stop()

	//HANDLING THE CLIENT MAKING THE RESERVATION
	fmt.Println("Starting client reciever")
	lis, err := net.Listen("tcp", ":3002") //React port + 2
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	gcrpServer := grpc.NewServer()
	//proto generates ReservationService Server
	proto.RegisterReservationServiceServer(gcrpServer, &server{})
	reflection.Register(gcrpServer)
	if err := gcrpServer.Serve(lis); err != nil {
		log.Fatalf("Error %v", err)
	}

	fmt.Println("Closing")
}

type server struct{}

func (s *server) MakeReservation(ctx context.Context, req *proto.BookTrip) (*proto.Trip, error) {
	log.Printf("Recieved request to book trip from client %s", req.PassengerName)
	booked_trip := proto.TripBooked{
		Trip: &proto.Trip{
			PassengerName: req.PassengerName,
			DriverName:    "Marek",
		},
	}
	jsonbytes, er := json.Marshal(booked_trip)
	if er != nil {
		log.Printf("Could not convert json %v", er)
		return &proto.Trip{}, errors.New("Could not convert JSON")
	}
	ps_ctx := context.Background()
	if _, er = topic.Publish(ps_ctx, &pubsub.Message{Data: jsonbytes}).Get(ps_ctx); er != nil {
		return &proto.Trip{}, errors.New("Could not push confirmation to Pub/Sub service")
	}
	return booked_trip.Trip, nil
}
