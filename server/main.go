package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/marchmiel/proto-playground/Config"
	"github.com/marchmiel/proto-playground/client/model"
	//"github.com/pkg/errors"
	//	"github.com/marchmiel/proto-playground/proto"
	"github.com/marchmiel/proto-playground/server/servgrpc"
	"github.com/marchmiel/proto-playground/server/servpubs"
	"github.com/marchmiel/proto-playground/server/servrest"
	"github.com/marchmiel/proto-playground/server/wrapper"
	"github.com/pkg/errors"
	//	"google.golang.org/grpc"
	//	"google.golang.org/grpc/reflection"
	"log"
	"time"
	//	"net"
	//	"net/http"
	//	"os"
	//"reflect"
)

var designatedDriver string = "Marek"

//var publisherMain message.Publisher

func main() {
	errorsChan := make(chan error, 3)
	collection := make(chan wrapper.ClientDataType, 100)

	go servgrpc.NewServePortChan("", collection, errorsChan)
	go servpubs.NewServePortChan("", collection, errorsChan)
	go servrest.NewServePortChan("", collection, errorsChan)

	stop, _ := context.WithTimeout(context.Background(), time.Second*2)
checkErrors:
	for {
		select {
		case <-stop.Done():
			break checkErrors
		case err := <-errorsChan:
			fmt.Println(errors.Wrap(err, " error in starting up endpoint"))
		}
	}
	log.Printf("server startups completed")

	logrus := watermill.NewStdLogger(false, false)
	publisherMain, err := googlecloud.NewPublisher(googlecloud.PublisherConfig{
		ProjectID: Config.PubSub_Project_Name,
	}, logrus)
	if err != nil {
		panic(err)
	}
	log.Printf("loop")
	ServeChannel(collection, publisherMain)
}

func ServeChannel(collection <-chan wrapper.ClientDataType, publisherMain message.Publisher) {
	driverID := 1
	for clientData := range collection {
		var bookTripRequest model.BookTripRequest
		err := clientData.Unload(&bookTripRequest)
		if err != nil {
			clientData.SendBack(errors.Wrap(err, "Could not unload BookTripRequest"))
			continue
		}
		log.Printf("Recieved request from %s", bookTripRequest.PassengerName)
		tripBookedResponse := model.TripBookedResponse{
			PassengerName: bookTripRequest.PassengerName,
			DriverName:    designatedDriver + string(driverID),
		}
		driverID += 1
		//PUBLISH MESSAGE - - - - - - -
		jsonbytes, err := json.Marshal(tripBookedResponse)
		if err != nil {
			clientData.SendBack(errors.Wrap(err, "Could not convert json for publisher"))
			continue
		}
		msg := message.NewMessage(clientData.CorrelationID(), jsonbytes)
		err = publisherMain.Publish(Config.Server_Publish_Topic, msg)
		if err != nil {
			clientData.SendBack(errors.Wrap(err, "Could not publish"))
			continue
		}
		//SEND DATA BACK TO CLIENT - - - - - -
		err = clientData.Load(&tripBookedResponse)
		if err != nil {
			log.Printf("Published a trip that could not be sent back")
			//clientData.SendBack(errors.Wrap(err, "Could not load TripBookedResponse"))
		}
	}
}

/*
func PublisherInit() {
	var err error
	if publisherMain, err = googlecloud.NewPublisher(googlecloud.PublisherConfig{
		ProjectID: Config.PubSub_Project_Name,
	}, watermill.NewStdLogger(false, false)); err != nil {
		panic(err)
	}
}

func HandleClientData(clientData wrapper.ClientDataType) (*model.TripBookedResponse, error) {
	var bookTripRequest model.BookTripRequest
	err := clientData.Unload(&bookTripRequest)
	if err != nil {
		return nil, err
	}
	tripBookedResponse := model.TripBookedResponse{
		PassengerName: bookTripRequest.PassengerName,
		DriverName:    designatedDriver,
	}
	err = clientData.Load(&tripBookedResponse)
	if err != nil {
		return nil, err
	}
	return &tripBookedResponse, nil
}

func Publish(tbr *model.TripBookedResponse, msg *message.Message) error { // in the form of (nil, msg) or (tbr, nil)
	var jm *message.Message
	jm = msg
	if jm == nil {
		jsonbytes, err := json.Marshal(tbr)
		if err != nil {
			return nil //, errors.Wrap(err, "Could not convert json")
		}
		jm = message.NewMessage(watermill.NewUUID(), jsonbytes)
	}
	return publisherMain.Publish(Config.Server_Publish_Topic, jm)
}

//NEEDS A BIT OF REFACTORING BELOW

func Grpc() {
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
	//var cli wrapper.ClientDataTyp
	s.l.Println("Recieveing gRPC request to book trip from " + req.PassengerName)
	cli := servgrpc.NewGrpcData(req)
	tbr, err := HandleClientData(cli)
	if err != nil {
		panic(err)
	}

	Publish(tbr, nil)
	return cli.GetResponse().(*proto.TripBooked), nil
}

func Pubs() {
	logger := watermill.NewStdLogger(false, false)
	subscriber, err := googlecloud.NewSubscriber(
		googlecloud.SubscriberConfig{
			GenerateSubscriptionName: func(topic string) string {
				return "server-puller"
			},
			ProjectID: Config.PubSub_Project_Name,
		},
		logger,
	)
	if err != nil {
		panic(err)
	}
	messages, err := subscriber.Subscribe(context.Background(), Config.Server_Pull_Topic)
	if err != nil {
		panic(err)
	}
	for m := range messages {
		//var cli wrapper.ClientDataTyp
		fmt.Println("Recieved m channel ")
		cli := servpubs.NewPubsData(m)
		_, err := HandleClientData(cli)
		if err != nil {
			panic(err)
		}
		//val := reflect.ValueOf(user).Elem()
		Publish(nil, cli.GetResponse().(*message.Message))
		log.Println("received message: %s, payload: %s", m.UUID, string(m.Payload))
		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		m.Ack()
	}
}

func Rest() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//var cli wrapper.ClientDataTyp
		cli := servrest.NewRestData(w, r)
		tbr, err := HandleClientData(cli)
		fmt.Println(err)
		if err != nil {
			panic(err)
		}
		Publish(tbr, nil)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(tbr)
	})
	http.ListenAndServe(":3003", nil)
}
*/
