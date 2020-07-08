package servpubs

import (
	"encoding/json"
	//"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	//	"github.com/marchmiel/proto-playground/client/model"
	"github.com/marchmiel/proto-playground/server/wrapper"
)

type pubsData struct {
	ReqMsg  *message.Message
	RespMsg *message.Message
}

func NewPubsData(req *message.Message) *pubsData {
	return &pubsData{ReqMsg: req}
}

func (p *pubsData) Unload(btr *wrapper.BookTripRequest) error {
	return json.Unmarshal(p.ReqMsg.Payload, btr)
}
func (p *pubsData) Load(tbr *wrapper.TripBookedResponse) error {
	jsonbytes, err := json.Marshal(tbr)
	p.RespMsg = message.NewMessage(p.ReqMsg.UUID, jsonbytes)
	return err
}
func (p *pubsData) Ret() interface{} {
	return p.RespMsg
}
