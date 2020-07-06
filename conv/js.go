package conv

import (
	"encoding/json"
	"errors"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

func JsonToMessage(class interface{}) (*message.Message, error) {
	jsonbytes, err := json.Marshal(class)
	if err != nil {
		return nil, errors.New("Could not convert json: " + err.Error())
	}
	msg := message.NewMessage(watermill.NewUUID(), jsonbytes)
	return msg, nil
}
