package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/majorsabbir/pubsub/pubsubpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

var reflectable bool = false
var tls bool = false

func loadEnv() {
	envErr := godotenv.Load("./.env")
	if envErr != nil {
		log.Fatalf("Error loading .env file %v", envErr)
	}
}

func parseConf(env string) bool {
	rawConf := os.Getenv(env)
	conf, err := strconv.ParseBool(rawConf)
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}
	return conf
}

func bootstrap() {
	// if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// load .env
	loadEnv()

	// parse reflection conf
	reflectable = parseConf("REFLECTION")

	// parse tls conf
	tls = parseConf("TLS")
}

type server struct{}

func (*server) PublishEvent(ctx context.Context, req *pubsubpb.PublishEventRequest) (*pubsubpb.PublishEventResponse, error) {
	fmt.Println("Trigger publish event")

	q := rdb.PubSubNumSub(context.Background(), req.Event.GetChannel())

	var streamable bool = false
	var subscriber_count int64 = 0

	for _, value := range q.Val() {
		if value > 0 {
			streamable = true
			subscriber_count = value
		}
	}

	if streamable {
		rdb.Publish(context.Background(), req.Event.GetChannel(), req.Event.GetMsg())
	}

	return &pubsubpb.PublishEventResponse{
		PublishEvent: &pubsubpb.PublishedEvent{
			Channel:         req.Event.GetChannel(),
			Msg:             req.Event.GetMsg(),
			SubscriberCount: subscriber_count,
		},
	}, nil
}

func receiveAndStream(sub *redis.PubSub, stream pubsubpb.PubsubService_ListenEventServer) {
	res, err := sub.ReceiveMessage(context.Background())
	if err != nil {
		fmt.Printf("Internal error: %v", err)
	}

	stream.Send(&pubsubpb.ListenEventResponse{
		Event: &pubsubpb.Event{
			Channel: res.Channel,
			Msg:     res.Payload,
		},
	})

	receiveAndStream(sub, stream)
}

func (*server) ListenEvent(req *pubsubpb.ListenEventRequest, stream pubsubpb.PubsubService_ListenEventServer) error {
	fmt.Println("Trigger listen event")

	sub := rdb.Subscribe(context.Background(), req.GetChannel())

	receiveAndStream(sub, stream)

	return nil
}

var rdb *redis.Client

func main() {
	bootstrap()

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}

	if tls {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatalf("Failed loading certificates: %v", sslErr)
			return
		}

		opts = append(opts, grpc.Creds(creds))
	}

	s := grpc.NewServer(opts...)
	pubsubpb.RegisterPubsubServiceServer(s, &server{})

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// reflection
	if reflectable {
		reflection.Register(s)
	}

	go func() {
		fmt.Println("Starting Server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch
	fmt.Println("\nStopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("End of Program")
}
