package core

import (
	"context"
	"log"
	"fmt"
	"github.com/fle4a/transaction-system/invoice-service/src/internal/types"
	"github.com/fle4a/transaction-system/invoice-service/src/configs"
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

func (db *DBPool) Invoice(txn types.InvoiceBody) error {
	conn, err := db.Pool.Acquire(db.Context)
	if err != nil {
		return err
	}
	defer conn.Release()
	_, err = conn.Exec(db.Context, "UPDATE wallets SET amount = amount + $1 WHERE wallet_id = $2 and currency = $3", txn.Amount, txn.ReceiverWalletId, txn.Currency)
	return err
}

func (db *DBPool) Close() {
	log.Println("Closing connect")
	db.Pool.Close()
}
