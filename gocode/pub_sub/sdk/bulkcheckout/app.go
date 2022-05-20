package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	dapr "github.com/dapr/go-sdk/client"
)

var (
	PUBSUB_NAME  = "orderpubsub"
	PUBSUB_TOPIC = "bulkorders"
)

func main() {
	// Wait few seconds for Dapr side car http server to be up
	time.Sleep(5 * time.Second)

	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()
	ctx := context.Background()
	for i := 1; i <= 1500; i += 2 {
		order := `{"orderId":` + strconv.Itoa(i) + `}`

		// Publish an event using Dapr pub/sub
		if err := client.PublishEvent(ctx, PUBSUB_NAME, PUBSUB_TOPIC, []byte(order)); err != nil {
			panic(err)
		}

		fmt.Println("Bulk Published data: ", order)

		time.Sleep(1 * time.Second)
	}
}
