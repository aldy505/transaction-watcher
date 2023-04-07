package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/twmb/franz-go/pkg/kgo"
	"log"
	"time"
)

type BalanceLog struct {
	CustomerNumber int64  `json:"customer_number"`
	Amount         int64  `json:"amount"`
	Timestamp      string `json:"timestamp"`
	TransactionID  int64  `json:"transaction_id"`
}

func produce(client *kgo.Client, balanceLog BalanceLog) error {
	payload, err := json.Marshal(balanceLog)
	if err != nil {
		return fmt.Errorf("marshaling json: %w", err)
	}

	ctx := context.Background()
	client.Produce(ctx, &kgo.Record{
		Value:     payload,
		Headers:   nil,
		Timestamp: time.Now(),
		Topic:     "balance",
		Context:   ctx,
	}, func(_ *kgo.Record, err error) {
		if err != nil {
			log.Println(err)
		}
	})

	return nil
}
