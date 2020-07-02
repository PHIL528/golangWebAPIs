package Config

var (
	//Publish/Pull are relative to server
	Server_Publish_Topic         = "events.TripBooked"
	Server_Pull_Topic            = "events.MakeReservation"
	GRPC_PORT             string = ":3002"          //Exposed
	Localhost_PubSub_PORT string = "localhost:8085" //Connecting to external
	PubSub_Project_Name          = "karhoo-local"
)
