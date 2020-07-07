package servRest

import (
	"github.com/marchmiel/proto-playground/client/model"
	"github.com/marchmiel/proto-playground/server/wrapper"
	"encoding/json"
    "fmt"
    "log"
    "net/http"
)

type restData {
	body *http.Request.Body

}
func (rest *restData) Unload(btr *model.BookTripRequest) error {
	return json.NewDecoder(rest.body).Decode(btr)
}
func (rest *restData) Load(tbr *model.TripBookedResponse) {
	
}
