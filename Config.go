package playground

var (
	publish_topic_name           = "events.TripBooked"
	pull_topic_name              = "events.MakeReservation"
	gRPC_PORT             string = ":3002"          //Exposed
	localhost_PubSub_PORT string = "localhost:8085" //Connecting to external
)
