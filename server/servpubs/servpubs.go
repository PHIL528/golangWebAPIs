package servpubs

import (
	"encoding/json"
	//"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/marchmiel/proto-playground/client/model"
	//"github.com/marchmiel/proto-playground/server/wrapper"
)

type pubsDataType struct {
	reqMsg  *message.Message
	respMsg *message.Message
	err     chan error
}

func NewPubsData(req *message.Message) *pubsDataType {
	return &pubsDataType{reqMsg: req, err: make(chan error)}
}

func (p *pubsDataType) Unload(btr *model.BookTripRequest) error {
	return json.Unmarshal(p.reqMsg.Payload, btr)
}
func (p *pubsDataType) Load(tbr *model.TripBookedResponse) error { //Not needed
	p.SendBack(nil)
	return nil
}
func (p *pubsDataType) CorrelationID() string { //Returns correlation ID of sender
	return p.reqMsg.UUID
}
func (p *pubsDataType) SendBack(err error) {
	p.err <- err
}

//jsonbytes, err := json.Marshal(tbr)
//p.respMsg = message.NewMessage(p.reqMsg.UUID, jsonbytes)
//return err
