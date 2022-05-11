package main

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/svid/jwtsvid"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

var (
	PUBSUB_NAME  = "order-pub-sub"
	PUBSUB_TOPIC = "orders"
)

func testSpiffe(addr string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	clientOptions := workloadapi.WithClientOptions(workloadapi.WithAddr(addr))

	fmt.Printf("CLIENT OPTIONS")

	jwtSource, err := workloadapi.NewJWTSource(
		ctx,
		clientOptions,
	)

	if err != nil {
		fmt.Printf("PANIC FROM NewJWTSource")
		panic(err)
	} else {
		fmt.Println("successfully got SPIFFE ID")
	}

	defer jwtSource.Close()

	serverID := spiffeid.RequireFromString("spiffe://iotedge/mqttbroker")
	audience := serverID.String()

	svid, err := jwtSource.FetchJWTSVID(ctx, jwtsvid.Params{
		Audience: audience,
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("message bus got SPIFFE ID - %s", svid.Marshal())
}

func main() {
	listener, err := net.Dial("unix", "/run/iotedge/sockets/workloadapi.sock")
	if err != nil {
		panic(err)
	}

	addr := fmt.Sprintf("%s://%s", listener.RemoteAddr().Network(), listener.RemoteAddr().String())
	fmt.Printf("%s", addr)

	testSpiffe(addr)

	// client, err := dapr.NewClient()
	// if err != nil {
	// 	panic(err)
	// }
	// defer client.Close()
	// ctx := context.Background()

	for i := 1; i <= 2000; i++ {
		order := `{"orderId":` + strconv.Itoa(i) + `}`

		// // Publish an event using Dapr pub/sub
		// if err := client.PublishEvent(ctx, PUBSUB_NAME, PUBSUB_TOPIC, []byte(order)); err != nil {
		// 	panic(err)
		// }

		fmt.Println("Published data: ", order)

		time.Sleep(2 * time.Second)
	}
}
