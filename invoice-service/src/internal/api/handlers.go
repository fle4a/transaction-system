package api

import (
	"encoding/json"
	"net/http"

	"github.com/fle4a/transaction-system/invoice-service/src/internal/core"
	"github.com/fle4a/transaction-system/invoice-service/src/internal/types"
)

func InvoiceHandler(db *core.DBPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var RequestData types.InvoiceBody
		err := json.NewDecoder(r.Body).Decode(&RequestData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		err = db.Invoice(RequestData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}
