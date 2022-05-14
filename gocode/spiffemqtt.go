package main

import (
	"context"
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/svid/jwtsvid"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

var (
	PUBSUB_NAME  = "order-pub-sub"
	PUBSUB_TOPIC = "orders"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func testSpiffe() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	serverID := spiffeid.RequireFromString("spiffe://iotedge/mqttbroker")

	svid, err := workloadapi.FetchJWTSVID(
		ctx,
		jwtsvid.Params{
			Audience: serverID.String(),
		},
		workloadapi.WithAddr("unix:///run/iotedge/sockets/workloadapi.sock"),
	)

	if err != nil {
		panic(err)
	} else {
		fmt.Println("MQTTE4K: successfully got SPIFFE ID")
	}

	fmt.Println("")
	fmt.Printf("message bus got SPIFFE ID - %s", svid.Marshal())
	fmt.Println("")
	fmt.Printf("SPIFFE ID: %s, AUDIENCE: %s, EXPIRY: %s", svid.ID, svid.Audience, svid.Expiry)
	fmt.Println("")

	opts := mqtt.NewClientOptions()
	opts.SetDefaultPublishHandler(f)
	opts.SetClientID(svid.ID.String())
	opts.SetCleanSession(true)
	opts.AddBroker("tcp://mqttbroker:1883")
	opts.SetUsername(svid.ID.String())
	opts.SetPassword(svid.Marshal())

	x := 0
	time_now := time.Now()
	expiry_time := svid.Expiry.UTC()

	for {
		client := mqtt.NewClient(opts)
		time_now = time.Now().UTC()

		if time_now == svid.Expiry.UTC() || time_now.After(expiry_time) {
			fmt.Println("TOKEN TIME IS EXPIRED")
		}

		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}

		// Subscribe to a topic
		if token := client.Subscribe("testtopic/#", 0, nil); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}

		// Publish a message
		message := fmt.Sprintf("MESSAGE# : %d", x)
		x = x + 1

		token := client.Publish("testtopic/1", 0, false, message)
		token.Wait()

		time.Sleep(250 * time.Millisecond)

		// Unscribe
		if token := client.Unsubscribe("testtopic/#"); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}

		// Disconnect
		client.Disconnect(250)
	}
}

func main() {
	testSpiffe()

}
