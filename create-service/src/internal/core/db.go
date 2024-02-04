package core

import (
	"context"
	"fmt"
	"log"

	"github.com/fle4a/transaction-system/create-service/src/configs"
	"github.com/fle4a/transaction-system/create-service/src/internal/types"
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

func (db *DBPool) createTransaction(txn types.Transaction) error {
	conn, err := db.Pool.Acquire(db.Context)
	if err != nil {
		return err
	}
	defer conn.Release()
	_, err = conn.Query(db.Context, "INSERT INTO transactions(transaction_id, sender_wallet_id, receiver_wallet_id, currency, amount, status) VALUES ($1, $2, $3, $4, $5, $6)", txn.ID, txn.SenderWalletId, txn.ReceiverWalletId, txn.Currency, txn.Amount, "Created")
	return err
}

func (db *DBPool) Close() {
	log.Println("Closing connect")
	db.Pool.Close()
}
