package types

import "github.com/google/uuid"

type BalanceRequest struct {
	WalletId uuid.UUID `json:"receiver_wallet_id"`
	Currency string    `json:"currency"`
}

type BalanceResponse struct {
	Actual float64 `json:"actual"`
	Frozen float64 `json:"frozen"`
}
