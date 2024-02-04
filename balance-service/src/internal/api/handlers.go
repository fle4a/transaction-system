package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fle4a/transaction-system/balance-service/src/internal/core"
	"github.com/fle4a/transaction-system/balance-service/src/internal/types"
)

func BalanceHandler(db *core.DBPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var RequestData types.BalanceRequest
		err := json.NewDecoder(r.Body).Decode(&RequestData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		log.Println(RequestData)
		data, err := db.GetBalance(RequestData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bytesData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write(bytesData)
	}
}
