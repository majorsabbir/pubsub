package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/majorsabbir/pubsub/pubsubpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Create blog function
func PublishEvent(c pubsubpb.PubsubServiceClient) {
	fmt.Println("Publishing an event")
	event := &pubsubpb.Event{
		Channel: "test",
		Msg:     "This is a test message",
	}

	res, err := c.PublishEvent(context.Background(), &pubsubpb.PublishEventRequest{
		Event: event,
	})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Event has been published: %v\n", res)
}

// List blog function
func ListenEvent(c pubsubpb.PubsubServiceClient) {
	fmt.Println("Listening event")

	res, err := c.ListenEvent(context.Background(), &pubsubpb.ListenEventRequest{
		Channel: "test",
	})
	if err != nil {
		fmt.Printf("Error happened while listening: %v\n", err)
	}
	for {
		msg, err := res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		fmt.Println(msg.GetEvent())
	}
}

func main() {
	fmt.Println("Hello from pubsub client")

	// load .env
	envErr := godotenv.Load("./.env")
	if envErr != nil {
		log.Fatalf("Error loading .env file %v", envErr)
	}

	// parse tls conf
	rawTlsConf := os.Getenv("TLS")
	parsedTlsConf, parseTlsErr := strconv.ParseBool(rawTlsConf)
	if envErr != parseTlsErr {
		log.Fatalf("Error loading .env file %v", parseTlsErr)
	}

	tls := parsedTlsConf
	opts := grpc.WithInsecure()

	if tls {
		credFile := "ssl/ca.crt" // Certificate Authority Trust Certificate
		creds, sslErr := credentials.NewClientTLSFromFile(credFile, "")
		if sslErr != nil {
			log.Fatalf("Error while loading CA trust certificates: %v", sslErr)
			return
		}

		opts = grpc.WithTransportCredentials(creds)
	}

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Failed to connect client: %v", err)
	}
	defer cc.Close()

	c := pubsubpb.NewPubsubServiceClient(cc)
	// ListenEvent(c)  // listen event
	PublishEvent(c) // publish event

}
