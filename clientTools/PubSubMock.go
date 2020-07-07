package clientTools

import (
	"github.com/ThreeDotsLabs/watermill/message"
)

type PubSubConnector interface {
	SetSubName(string)
	GetSubName() string
	SendReservation(*message.Message) (<-chan *message.Message, error)

	SetConnType(string)
	GetConnType() string
}

//counterfeiter:generate ./../.. PubSubConnector
//go run github.com/maxbrunsfeld/counterfeiter/v6 ./clientTools/PubSubMock.go PubSubConnector
