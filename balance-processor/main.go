package main

import (
	"context"
	"errors"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kversion"
	"log"
	"os"
	"os/signal"
	"strings"
)

func main() {
	kafkaHost, ok := os.LookupEnv("KAFKA_HOST")
	if !ok {
		kafkaHost = "localhost:9092"
	}

	kafkaOptions := []kgo.Opt{
		kgo.MinVersions(kversion.V0_11_0()),
		kgo.SeedBrokers(strings.Split(kafkaHost, ",")...),
		kgo.ConsumeTopics("transactions"),
	}

	kafkaClient, err := kgo.NewClient(kafkaOptions...)
	if err != nil {
		log.Fatalf("initiating kafka client: %s", err.Error())
	}
	defer kafkaClient.Close()

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, os.Interrupt)

	go func() {
		<-exitSignal

		kafkaClient.Close()
	}()

	log.Println("Ready, waiting for messages...")

	for {
		ctx := context.Background()
		fetches := kafkaClient.PollRecords(ctx, 128)

		switch {
		case fetches.IsClientClosed():
			break
		case fetches.Err() != nil:
			fetches.EachError(func(topic string, partition int32, err error) {
				if !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
					log.Printf("Error: topic: %s, partition: %d, message: %v", topic, partition, err)
				}
			})
			continue
		}

		for _, msg := range fetches.Records() {
			err := consume(kafkaClient, msg.Value)
			if err != nil {
				log.Println(err)
				continue
			}

			log.Printf("Processed: %s", string(msg.Value))
		}
	}
}
