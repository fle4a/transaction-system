package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fle4a/transaction-system/withdraw-service/src/configs"
	"github.com/fle4a/transaction-system/withdraw-service/src/internal/types"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DBPool struct {
	Pool    *pgxpool.Pool
	Context context.Context
	Dburl   string
}

func CreateDBURL(config *configs.Config) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		config.Database.User,
		config.Database.Pass,
		config.Database.Host,
		config.Database.Port,
		config.Database.Db)
}

func NewPool(context context.Context, dburl string) *DBPool {
	return &DBPool{
		Context: context,
		Dburl:   dburl,
	}
}

func (db *DBPool) Init() error {
	pool, err := pgxpool.Connect(db.Context, db.Dburl)
	if err != nil {
		return err
	}
	db.Pool = pool

	return nil
}

func (db *DBPool) ProcessTransaction(txn types.Transaction) error {

	conn, err := db.Pool.Acquire(db.Context)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(db.Context)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(db.Context)
		} else {
			tx.Commit(db.Context)
		}
	}()

	var senderBalance float64
	if err := tx.QueryRow(db.Context, "SELECT amount FROM wallets WHERE wallet_id = $1 AND currency = $2", txn.SenderWalletId, txn.Currency).Scan(&senderBalance); err != nil {
		return err
	}

	if senderBalance < txn.Amount {
		return fmt.Errorf("Недостаточно средств")
	}

	_, err = tx.Exec(db.Context, "UPDATE wallets SET amount = amount - $1 WHERE wallet_id = $2", txn.Amount, txn.SenderWalletId)
	if err != nil {
		return err
	}

	body := types.InvoiceBody{
		ReceiverWalletId: txn.ReceiverWalletId,
		Currency:         txn.Currency,
		Amount:           txn.Amount,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://invoice-service:8002/invoice", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 202 {
		return fmt.Errorf("Error from invoice service")
	}

	return nil
}

func (db *DBPool) ChangeStatus(txn types.Transaction, status string) error {
	conn, err := db.Pool.Acquire(db.Context)
	if err != nil {
		log.Printf("Error get connection: %v", err)
		return err
	}
	defer conn.Release()
	_, err = conn.Query(db.Context, `UPDATE transactions
									SET status = $1
									WHERE sender_wallet_id = $2 and
									receiver_wallet_id = $3 and
									currency = $4`,
		status, txn.SenderWalletId, txn.ReceiverWalletId, txn.Currency)
	return err
}

func (db *DBPool) Close() {
	log.Println("Closing connect")
	db.Pool.Close()
}
