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
	var err error
	publish_topic, err = setupPublisher() //To publish confirmed requests, to be recieved by client + listener
	if err != nil {
		panic(err)
	}
	pull_subscription, err := setupPuller() //To recieve requests from client
	if err != nil {
		panic(err)
	}

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

func setupPublisher() (*pubsub.Topic, error) {
	c_log := log.New(os.Stdout, "setupPublisher", log.LstdFlags)
	os.Setenv("PUBSUB_EMULATOR_HOST", Config.Localhost_PubSub_PORT)

	ps_ctx := context.Background()
	notified := false
	var client *pubsub.Client
	var err error
	for { //I put the for loop in hear to loop until the PubSub is started, but its useless because err returns nil regardless of whether or not the PubSub terminal is started
		client, err = pubsub.NewClient(ps_ctx, "karhoo-local")
		if err == nil {
			break
		} else if !notified {
			c_log.Printf("setupPublisher: Failed to create pubsub client, %v", err)
			c_log.Printf("setupPublisher Perhaps the PubSub terminal has not yet been started, will reattempt conncetion once per second")
			notified = true
		}
		time.Sleep(time.Second)
	}
	topic, err := client.CreateTopic(ps_ctx, Config.Server_Publish_Topic)
	if err != nil {

	}
	var return_err error = nil
	if err != nil {
		c_log.Printf("Failed to create topic %v", err)
		c_log.Printf("Perhaps the topic already exists, joining topic instead of creating")
		topic = client.Topic(Config.Server_Publish_Topic)
		if exists, err := topic.Exists(ps_ctx); !exists {
			return_err = errors.New("Cannot create topic and topic does not exist, " + err.Error())
		}
	}
	return topic, return_err
}

func setupPuller() (*pubsub.Subscription, error) {
	fmt.Println("")
	sub, _, err := Config.GetSubscriptionToTopic(context.Background(), Config.Server_Pull_Topic, "server-pull", true)
	if err != nil {
		panic(err)
	}
	return sub, err
}
