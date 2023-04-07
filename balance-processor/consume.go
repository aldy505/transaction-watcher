package main

import (
	"encoding/json"
	"fmt"
	"github.com/twmb/franz-go/pkg/kgo"
	"strings"
	"time"
)

type TransactionLog struct {
	TransactionID     int64  `json:"transaction_id"`
	TransactionType   string `json:"transaction_type"`
	CustomerNumber    int64  `json:"customer_number"`
	TransactionAmount int64  `json:"transaction_amount"`
	Timestamp         string `json:"timestamp"`
}

func consume(client *kgo.Client, payload []byte) error {
	var transactionLog TransactionLog
	if err := json.Unmarshal(payload, &transactionLog); err != nil {
		return fmt.Errorf("invalid payload: %s", err.Error())
	}

	var amount int64
	switch strings.ToUpper(transactionLog.TransactionType) {
	case "TOP_UP":
		amount = transactionLog.TransactionAmount
		break
	case "TRANSFER":
		fallthrough
	case "WITHDRAW":
		fallthrough
	case "FEE":
		amount = -1 * transactionLog.TransactionAmount
		break
	default:
		return fmt.Errorf("unknown transaction type: %s", transactionLog.TransactionType)
	}

	return produce(client, BalanceLog{
		CustomerNumber: transactionLog.CustomerNumber,
		Amount:         amount,
		Timestamp:      time.Now().Format(time.RFC3339),
		TransactionID:  transactionLog.TransactionID,
	})
}
