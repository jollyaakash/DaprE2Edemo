package main

import (
	"context"
	"log"
	"net/http"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
)

var sub = &common.Subscription{
	PubsubName: "orderpubsub",
	Topic:      "orders",
	Route:      "orders",
}

func main() {
	s := daprd.NewService(":6001")
	if err := s.AddTopicEventHandler(sub, eventHandler); err != nil {
		log.Fatalf("error adding topic subscription: %v", err)
	}
	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error listenning: %v", err)
	}
}

func eventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	log.Printf("Subscriber received message PubSubName: %s, Topic: %s, Message: %s", e.PubsubName, e.Topic, e.Data)
	return false, nil
}
