package core

import (
	"context"
	"fmt"
	"log"

	"github.com/fle4a/transaction-system/balance-service/src/configs"
	"github.com/fle4a/transaction-system/balance-service/src/internal/types"
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

func (db *DBPool) GetBalance(data types.BalanceRequest) (types.BalanceResponse, error) {
	var result types.BalanceResponse
	conn, err := db.Pool.Acquire(db.Context)
	if err != nil {
		return result, err
	}
	defer conn.Release()
	err = conn.QueryRow(db.Context, "SELECT amount FROM wallets WHERE wallet_id = $1 and currency = $2;", data.WalletId, data.Currency).Scan(&result.Actual)
	if err != nil {
		return result, err
	}
	err = conn.QueryRow(db.Context, `
	SELECT 
		COALESCE(SUM(CASE
        	WHEN receiver_wallet_id = $1 THEN amount
        	WHEN sender_wallet_id = $1 THEN -amount
        	ELSE 0
    	END), 0.0) AS user_balance
		FROM transactions
	WHERE status = 'Created' 
		AND currency = $2
		AND (receiver_wallet_id = $1 OR sender_wallet_id = $1);
	`, data.WalletId, data.Currency).Scan(&result.Frozen)
	if err != nil {
		return result, err
	}
	return result, err
}

func (db *DBPool) Close() {
	log.Println("Closing connect")
	db.Pool.Close()
}
