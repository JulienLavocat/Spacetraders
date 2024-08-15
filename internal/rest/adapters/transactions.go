package adapters

import (
	"time"

	"github.com/julienlavocat/spacetraders/.gen/spacetraders/public/model"
	"github.com/julienlavocat/spacetraders/internal/api"
)

type TransactionsList struct {
	Transactions []Transaction `json:"transactions"`
	Revenue      int64         `json:"revenue"`
	Expenses     int64         `json:"expenses"`
	FuelExpenses int64         `json:"fuelExpenses"`
}

type Transaction struct {
	Timestamp     time.Time `json:"timestamp"`
	CorrelationId *string   `json:"correlationId"`
	Waypoint      string    `json:"waypoint"`
	Product       string    `json:"product"`
	Type          string    `json:"type"`
	Ship          string    `json:"ship"`
	AgentBalance  int64     `json:"agentBalance"`
	Id            int32     `json:"id"`
	TotalPrice    int32     `json:"totalPrice"`
	Amount        int32     `json:"amount"`
	PricePerUnit  int32     `json:"pricePerUnit"`
}

func AdaptTransaction(tx model.Transactions) Transaction {
	return Transaction{
		Timestamp:     tx.Timestamp,
		Waypoint:      tx.Waypoint,
		Product:       tx.Product,
		Type:          tx.Type,
		Ship:          tx.Ship,
		CorrelationId: tx.CorrelationID,
		Id:            tx.ID,
		TotalPrice:    tx.TotalPrice,
		AgentBalance:  tx.AgentBalance,
		Amount:        tx.Amount,
		PricePerUnit:  tx.PricePerUnit,
	}
}

func AdaptTransactions(txs []model.Transactions) TransactionsList {
	revenue := int64(0)
	expenses := int64(0)
	fuelExpenses := int64(0)

	transactions := make([]Transaction, len(txs))
	for i := range txs {
		transaction := txs[i]
		transactions[i] = AdaptTransaction(txs[i])
		if transaction.Type == "PURCHASE" {
			expenses += int64(transaction.TotalPrice)
			if transaction.Product == string(api.FUEL) {
				fuelExpenses += int64(transaction.TotalPrice)
			}
		}

		if transaction.Type == "SELL" {
			revenue += int64(transaction.TotalPrice)
		}
	}

	return TransactionsList{
		Transactions: transactions,
		Revenue:      revenue,
		Expenses:     expenses,
		FuelExpenses: fuelExpenses,
	}
}
