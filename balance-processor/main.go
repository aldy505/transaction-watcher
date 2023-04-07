package main

import (
	"context"
	"errors"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kversion"
	"log"
	"os"
	"os/signal"
	"strings"
)

func main() {
	kafkaHost, ok := os.LookupEnv("KAFKA_ADDRESSES")
	if !ok {
		kafkaHost = "localhost:9092"
	}

	kafkaOptions := []kgo.Opt{
		kgo.MinVersions(kversion.V0_11_0()),
		kgo.SeedBrokers(kafkaHost),
		kgo.ConsumeTopics("transactions"),
		kgo.AllowAutoTopicCreation(),
		kgo.WithLogger(kgo.BasicLogger(os.Stdout, kgo.LogLevelWarn, nil)),
		kgo.ClientID("balance-processor"),
	}

	kafkaClient, err := kgo.NewClient(kafkaOptions...)
	if err != nil {
		log.Fatalf("initiating kafka client: %s", err.Error())
	}
	defer kafkaClient.Close()

	// Create balance topic
	kafkaAdmin := kadm.NewClient(kafkaClient)
	_, err = kafkaAdmin.CreateTopic(context.Background(), 1, 1, nil, "balance")
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		kafkaClient.Close()

		log.Fatalf("Creating 'balance' topic: %s", err.Error())
	}

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
