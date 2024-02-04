package types

import (
	"github.com/google/uuid"
)

type Transaction struct {
	ID               uuid.UUID `json:"transaction_id"`
	SenderWalletId   uuid.UUID `json:"sender_wallet_id"`
	ReceiverWalletId uuid.UUID `json:"receiver_wallet_id"`
	Currency         string    `json:"currency"`
	Amount           float64   `json:"amount"`
}
