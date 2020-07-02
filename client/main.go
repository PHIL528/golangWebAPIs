package main

import (
	"cloud.google.com/go/pubsub"
	"encoding/json"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"proto-playground/Config"
	"proto-playground/proto"
	"strings"
)

func main() {
	route := os.Args[1]
	client_fname := strings.ToLower(os.Args[2])
	if route == "grpc" {
		send_via_gRPC(client_fname)
	} else if route == "pubsub" {
		send_via_PubSub(client_fname)
	} else {
		log.Fatalf(route, " is neither gRPC nor PubSub")
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
		return nil, errors.New("Error making new client "+ err.Error()))
	}
	send_topic := client.Topic(Config.Server_Pull_Topic)
	exist, err := topic.Exists(ps_ctx)
	if err != nil {
		return nil, errors.New("Error in checking if topic exists %v", err)
	} else if !exist {
		return nil, errors.New("events.MakeReservation topic has not been created")
	}
	defer send_topic.Stop()
	request := proto.BookTrip{
		PassengerName: client_name,
	}
	jsonbytes, err := json.Marshal(request)
	if err != nil {
		log.Printf("Could not convert json %v", err)
		return 
	}
	//ctx := context.Background() //Not sure if I should use a new context here.

	recieve_topic := client.Topic(ps_ctx, Config.Server_Pull_Topic)
	defer recieve_topic.Stop()
	if exists, err := topic.Exists(ps_ctx); !exists {
		return nil, errors.New("Reciever topic does not exist %v", err)
	}
	if err != nil {
		return nil, errors.New("Other error with topic checking")
	}
	sub := client.Subscription(client_name)
	exists, err := sub.Exists(ps_ctx)
	if err != nil {
		return nil, errors.New("topic checking error")
	}
	if !exists {
		fmt.Println("Creating new listener")
		sub, err := client.CreateSubscription(ps_ctx, client_name, pubsub.SubscriptionConfig{Topic: recieve_topic})
		if err != nil {
			return nil, errors.New("cannot create subscription to recieve message")
		}
	}
	

	if _, err = send_topic.Publish(ps_ctx, &pubsub.Message{Data: jsonbytes}).Get(ctx); err != nil {
		return nil, errors.New("Publishing errors "+err)
	}

	var trip proto.TripBooked 
	var recievedConfirmation bool 

	var i int = 0
	for {                                      //This would be a great place to implement WithTimeout()
		time.Sleep(time.Second)
		
		ctxx := context.Background()
	//	client, _ := pubsub.NewClient(ctx, "karhoo-local")
		err := s.Receive(ctxx, func(ctxxx context.Context, msg *pubsub.Message) {
			var TripBooked proto.TripBooked
			er := json.Unmarshal([]byte(msg.Data), &TripBooked)
			if er != nil {
				log.Printf("Could not unmarshal JSON, cannot confirm trip %v", er)
			} else {
				if TripBooked.PassengerName == client_name{
					trip = TripBooked
					recievedConfirmation = true
				}
				msg.Ack()
			}
		})
		if recievedConfirmation {
			break
		}
		i = i + 1
		if i>10 {
			return nil, errors.New("Cancel by timeout")
		}
	}
	return trip, nil
}
