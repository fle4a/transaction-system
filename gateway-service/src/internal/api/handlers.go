package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fle4a/transaction-system/gateway-service/src/internal/core"
	"github.com/fle4a/transaction-system/gateway-service/src/internal/types"
	"github.com/google/uuid"
)

func TransactionHandler(p *core.KafkaProducer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var txn types.Transaction
		err := json.NewDecoder(r.Body).Decode(&txn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		txn.ID = uuid.New()
		if txn.SenderWalletId == uuid.Nil || txn.ReceiverWalletId == uuid.Nil || txn.Currency == "" || txn.Amount <= 0 {
			http.Error(w, "Missing or invalid parameters", http.StatusBadRequest)
			return
		}
		err = p.Produce(txn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("Send to kafka: %s", txn.ID)
		w.WriteHeader(http.StatusAccepted)
		response := map[string]string{
			"message": "Transaction is being processed",
			"id":      txn.ID.String(),
		}
		responseJSON, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
	}
}

func BalanceHandler(w http.ResponseWriter, r *http.Request) {
	var blf types.BalanceForm
	err := json.NewDecoder(r.Body).Decode(&blf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	blfBytes, err := json.Marshal(blf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp, err := http.Post("http://balance-service:8003/balance", "application/json", bytes.NewBuffer(blfBytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}
