package Config

import (
	"cloud.google.com/go/pubsub"
	"context"
	"errors"
	"fmt"
)

var (
	//Publish/Pull are relative to server
	Server_Publish_Topic         = "events.TripBooked"
	Server_Pull_Topic            = "events.MakeReservation"
	GRPC_PORT             string = ":3002"          //Exposed
	Localhost_PubSub_PORT string = "localhost:8085" //Connecting to external
	PubSub_Project_Name          = "karhoo-local"
)

func GetTopic(ctx context.Context, top string, admin bool) (*pubsub.Topic, context.Context, error) {
	client, err := pubsub.NewClient(ctx, top)
	if err != nil {
		fmt.Println("1")
		return nil, nil, errors.New("Unable to setup connection with PubSub")
	}
	topic := client.Topic(top)
	exists, err := topic.Exists(ctx)
	if err != nil {
		fmt.Println("2")
		return nil, nil, errors.New("Could not check if topic exists")
	} else if !exists {
		if admin {
			topic, err = client.CreateTopic(ctx, top)
			if err != nil {
				fmt.Println("3")
				return nil, nil, errors.New("Could not create topic")
			}
		} else {
			fmt.Println("4")
			return nil, nil, errors.New("Topic has not been created yet")
		}
	}
	return topic, ctx, nil
}

func GetSubscriptionToTopic(ctx context.Context, top string, subID string, admin bool) (*pubsub.Subscription, context.Context, error) {
	client, err := pubsub.NewClient(ctx, PubSub_Project_Name)
	if err != nil {
		return nil, nil, errors.New("Could not up client")
	}
	topic := client.Topic(top)
	exists, err := topic.Exists(ctx)
	if err != nil {
		return nil, nil, errors.New("Other error with topic checking")
	}
	if !exists {
		if admin {
			topic, err = client.CreateTopic(ctx, top)
			if err != nil {
				return nil, nil, err
			}
		} else {
			return nil, nil, errors.New("Topic does not exist")
		}
	}
	sub := client.Subscription(subID)
	exists, err = sub.Exists(ctx)
	if err != nil {
		return nil, nil, errors.New("Error checking sub")
	}
	if !exists {
		sub, err = client.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{Topic: topic})
		if err != nil {
			return nil, nil, errors.New("Error creating subscription")
		}
	}
	return sub, ctx, nil
}
