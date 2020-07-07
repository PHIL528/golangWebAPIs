package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/marchmiel/proto-playground/Config"
	"github.com/marchmiel/proto-playground/clientTools"
	"log"
	"os"
	"strings"
	"time"
	//	"github.com/marchmiel/proto-playground/client/clientTools"
	"github.com/marchmiel/proto-playground/conv"
	"github.com/marchmiel/proto-playground/proto"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("client")
	route := strings.ToLower(os.Args[1])
	var err error
	var trip *proto.TripBooked
	fmt.Println("preroute")
	if route == "grpc" {
		client_fname := strings.ToLower(os.Args[2])
		trip, err = send_via_gRPC(client_fname)
	} else if route == "pubsub" {
		client_fname := strings.ToLower(os.Args[2])
		fmt.Println(client_fname)
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
	res, con, err := VarClientMaker.MakeClient()
	if err != nil {
		panic(err)
	}
	defer con.Close()

	BookTripRequest := proto.BookTrip{
		PassengerName: client_name,
	}
	TripBooked, err := res.MakeReservation(context.Background(), &BookTripRequest)
	if err != nil {
		panic(err)
	}
	return TripBooked, nil
}

var VarClientMaker clientTools.ClientMaker = NewClientMaker()

type clientMaker struct {
	//Conn        *grpc.ClientConn
	//ResClient   proto.ReservationServiceClient
	HandlerType string
}

func NewClientMaker() clientTools.ClientMaker {
	return &clientMaker{HandlerType: "Operational"}
}
func (g *clientMaker) MakeClient() (proto.ReservationServiceClient, *grpc.ClientConn, error) {
	var con *grpc.ClientConn
	con, err := grpc.Dial("localhost"+Config.GRPC_PORT, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	//g.Conn = con
	res := proto.NewReservationServiceClient(con)
	return res, con, err
}

func send_via_PubSub(client_name string) (*proto.TripBooked, error) {
	fmt.Println("STARTING PUBSUB")
	os.Setenv("PUBSUB_EMULATOR_HOST", Config.Localhost_PubSub_PORT)

	bookTripRequest := proto.BookTrip{
		PassengerName: client_name,
	}
	msg, err := conv.ToJsonToMessage(bookTripRequest)
	if err != nil {
		return nil, err
	}
	fmt.Println(VarPubSubConnector.GetConnType())
	VarPubSubConnector.SetSubName(client_name)
	messages, err := VarPubSubConnector.SendReservation(msg)

	var TripBooked proto.TripBooked
	for msg := range messages {
		err := json.Unmarshal(msg.Payload, &TripBooked)
		if err != nil {
			return nil, errors.New("Failed to unmarshal JSON")
		} else if TripBooked.Trip.PassengerName == client_name {
			log.Printf("Recieved trip confirmation, driver is " + TripBooked.Trip.DriverName)
			msg.Ack()
			break
		} else {
			log.Printf("Trip was from different client")
		}
		msg.Ack()
	}
	return &TripBooked, nil
}

var VarPubSubConnector clientTools.PubSubConnector = NewPubSubConnector()

type pubSubConnector struct {
	subscriptionName string
	ConnType         string
}

func NewPubSubConnector() clientTools.PubSubConnector {
	return &pubSubConnector{ConnType: "Operational"}
}
func (p *pubSubConnector) SetSubName(n string) {
	p.subscriptionName = n
}
func (p *pubSubConnector) GetSubName() string {
	return p.subscriptionName
}
func (p pubSubConnector) SetConnType(t string) {
	p.ConnType = t
}
func (p *pubSubConnector) GetConnType() string {
	return p.ConnType
}
func (p *pubSubConnector) SendReservation(msg *message.Message) (<-chan *message.Message, error) {
	logger := watermill.NewStdLogger(false, false)
	subscriber, err := googlecloud.NewSubscriber(
		googlecloud.SubscriberConfig{
			GenerateSubscriptionName: func(topic string) string {
				return p.subscriptionName
			},
			ProjectID: Config.PubSub_Project_Name,
		},
		logger,
	)
	if err != nil {
		return nil, errors.New("Could not create subscriber")
	}
	messages, err := subscriber.Subscribe(context.Background(), Config.Server_Publish_Topic)
	if err != nil {
		return nil, errors.New("Could not create subscriber message channel")
	}
	//PUBLISHER
	publisher, err := googlecloud.NewPublisher(googlecloud.PublisherConfig{
		ProjectID: Config.PubSub_Project_Name,
	}, logger)
	if err != nil {
		return nil, errors.New("Could not create publisher")
	}
	if err := publisher.Publish(Config.Server_Pull_Topic, msg); err != nil {
		return nil, errors.New("Could not push request to Pub/Sub service")
	}
	return messages, nil
}
