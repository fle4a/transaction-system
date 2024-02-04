package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/fle4a/transaction-system/tests/src/configs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Wallet struct {
	ID       uuid.UUID
	Currency string
	Amount   float64
}

type WithdrawData struct {
	SenderWalletId   uuid.UUID `json:"sender_wallet_id"`
	ReceiverWalletId uuid.UUID `json:"receiver_wallet_id"`
	Currency         string    `json:"currency"`
	Amount           float64   `json:"amount"`
}

type Balance struct {
	WalletID uuid.UUID `json:"receiver_wallet_id"`
	Currency string    `json:"currency"`
}

type BalanceResponse struct {
	Actual float64 `json:"actual"`
	Frozen float64 `json:"frozen"`
}

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

func (db *DBPool) Close() {
	log.Println("Closing connect")
	db.Pool.Close()
}

var db *DBPool

func main() {
	config, err := configs.ReadConfig()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	db = NewPool(context.Background(), CreateDBURL(config))
	if err := db.Init(); err != nil {
		log.Println(err)
		panic(err)
	}
	defer db.Close()

	if err := db.createTables(); err != nil {
		log.Fatalf("Error create tables %v\n", err)
	}

	const lenght = 10
	min := 30.
	max := 200.
	var Wallets [lenght]Wallet
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < lenght; i++ {
		wallet := Wallet{
			ID:       uuid.New(),
			Currency: "EUR",
			Amount:   min + rand.Float64()*(max-min),
		}
		Wallets[i] = wallet
		fmt.Printf("Wallet: %s\n", wallet.ID)
		err = db.WalletAdd(wallet)
		if err != nil {
			fmt.Printf("Error add wallet: %v\n", err)
		}
	}
	cnt := 0
	numTransactions := 1000
	for i := 0; i < numTransactions; i++ {
		sender := rand.Intn(lenght)
		receiver := rand.Intn(lenght)
		if sender == receiver {
			continue
		}

		senderWallet := Wallets[sender].ID
		receiverWallet := Wallets[receiver].ID
		currency := "EUR"
		amount := min + rand.Float64()*(max-min)
		body := WithdrawData{
			SenderWalletId:   senderWallet,
			ReceiverWalletId: receiverWallet,
			Currency:         currency,
			Amount:           amount,
		}
		cnt = cnt + 1
		err = SendTransaction(body)
		if err != nil {
			log.Println(err)
		}
	}
	for i := 0; i < lenght; i++ {
		data := Balance{
			WalletID: Wallets[i].ID,
			Currency: "EUR",
		}
		err = GetBalance(data)
		if err != nil {
			log.Println(err)
		}
	}
}

func (db *DBPool) createTables() error {
	conn, err := db.Pool.Acquire(db.Context)
	if err != nil {
		return err
	}
	defer conn.Release()
	_, err = conn.Exec(db.Context, `
	CREATE TABLE IF NOT EXISTS wallets (
		wallet_id UUID NOT NULL,
		currency  VARCHAR(255) NOT NULL,
		amount    NUMERIC(18, 2) NOT NULL,
		PRIMARY KEY (wallet_id, currency)
	);

	CREATE TABLE IF NOT EXISTS transactions (
		transaction_id     UUID NOT NULL,
		sender_wallet_id   UUID NOT NULL,
		receiver_wallet_id UUID NOT NULL,
		currency           VARCHAR(255) NOT NULL,
		amount             NUMERIC(18, 2) NOT NULL,
		status             VARCHAR(10) NOT NULL,
		PRIMARY KEY (transaction_id),
		FOREIGN KEY (sender_wallet_id, currency) REFERENCES wallets(wallet_id, currency),
		FOREIGN KEY (receiver_wallet_id, currency) REFERENCES wallets(wallet_id, currency)
	);

	CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
	`)
	return err
}

func (db *DBPool) WalletAdd(data Wallet) error {
	conn, err := db.Pool.Acquire(db.Context)
	if err != nil {
		return err
	}
	defer conn.Release()
	_, err = conn.Exec(db.Context,
		`INSERT INTO wallets(wallet_id, currency, amount)
		values ($1, $2, $3)`, data.ID, data.Currency, data.Amount)
	return err
}

func SendTransaction(data WithdrawData) error {
	bytesBody, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error parsing: %v\n", err)
	}

	resp, err := http.Post("http://localhost:8080/withdraw", "application/json", bytes.NewBuffer(bytesBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 202 {
		log.Printf("ERROR WITHDRAW %s\n", resp.Status)
	}
	return err
}

func GetBalance(data Balance) error {
	bytesBody, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error parsing: %v\n", err)
	}

	resp, err := http.Post("http://localhost:8080/balance", "application/json", bytes.NewBuffer(bytesBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 202 {
		log.Printf("ERROR BALANCE %s\n", resp.Status)
		data, err := io.ReadAll(resp.Body)
		log.Printf("Response body %s", data)
		return err
	}
	var respData BalanceResponse
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return err
	}
	log.Printf("Balance for wallet %s and currency %s:\n%v", data.WalletID, data.Currency, respData)
	return err
}
