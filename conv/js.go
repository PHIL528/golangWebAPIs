package conv

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

func ToJsonToMessage(emptyset interface{}) (*message.Message, error) {
	jsonbytes, err := json.Marshal(emptyset)
	if err != nil {
		return nil, errors.New("Could not convert json: " + err.Error())
	}
	msg := message.NewMessage(watermill.NewUUID(), jsonbytes)
	return msg, nil
}

func RecoverFromMessage(msg *message.Message, emptyset interface{}) (*interface{}, error) {
	err := json.Unmarshal(msg.Payload, &emptyset)
	if err != nil {
		fmt.Println(err)
		fmt.Println("READ ABOVE")
		return nil, errors.New("Failed to unmarshal JSON")
	}
	return &emptyset, nil
}
