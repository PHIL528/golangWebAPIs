package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"proto-playground/Config"
	"proto-playground/proto"
	"strings"
	"time"

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
	//ctx := context.Background() //Not sure if I should use a new context here.

	recieve_topic := client.Topic(Config.Server_Pull_Topic)
	if exists, err := recieve_topic.Exists(ps_ctx); !exists {
		return nil, errors.New("Reciever topic does not exist " + err.Error())
	}
	if err != nil {
		return nil, errors.New("Other error with topic checking" + err.Error())
	}
	nctx := context.Background()
	sub := client.Subscription(client_name)
	exists, err := sub.Exists(nctx)
	if err != nil {
		return nil, errors.New("topic checking error" + err.Error())
	}
	if !exists {
		fmt.Println("Creating new listener")
		sub, err = client.CreateSubscription(nctx, client_name, pubsub.SubscriptionConfig{Topic: recieve_topic})
		if err != nil {
			return nil, errors.New("cannot create subscription to recieve message" + err.Error())
		}
	}

	if _, err = send_topic.Publish(ps_ctx, &pubsub.Message{Data: jsonbytes}).Get(ps_ctx); err != nil {
		return nil, errors.New("Publishing errors " + err.Error())
	}
	time.Sleep(time.Second * 10)
	trip := proto.TripBooked{}
	var recievedConfirmation bool

	var i int = 0
	for { //This would be a great place to implement WithTimeout()
		time.Sleep(time.Second)

		ctxx := context.Background()
		//	client, _ := pubsub.NewClient(ctx, "karhoo-local")
		fmt.Println("THis is the sub ")
		fmt.Println(sub)
		err := sub.Receive(ctxx, func(ctxxx context.Context, msg *pubsub.Message) {
			thus := &proto.TripBooked{}
			er := json.Unmarshal([]byte(msg.Data), thus)
			if er != nil {
				log.Printf("Could not unmarshal JSON, cannot confirm trip %v", er)
			} else if thus == nil {
				//hi
			} else {
				fmt.Println("THIS IS TRIP BOOKED")
				fmt.Println(*thus)
				if thus.Trip.PassengerName == client_name {
					trip.Trip.PassengerName = thus.Trip.PassengerName
					trip.Trip.DriverName = thus.Trip.DriverName
					recievedConfirmation = true
				}
				msg.Ack()
			}
		})
		if err == err {

		}
		if recievedConfirmation {
			break
		}
		i = i + 1
		if i > 10 {
			return nil, errors.New("Cancel by timeout")
		}
	}
	return &trip, nil
}
