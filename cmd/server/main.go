package main

import (
	"database/sql"
	"fmt"
	"os"
	"subscriptionplus/server/infra/logger"
	"subscriptionplus/server/infra/store/postgres"
	"subscriptionplus/server/infra/store/postgres/store"

	"github.com/joho/godotenv"
)

var DefaultPortStr string = "8080"

type httpServer struct {
	db     *sql.DB
	logger *logger.Logger
	store  store.Storage
}

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// logger
	log := logger.NewLogger()

	// connection db
	db, err := postgres.New(os.Getenv("DB_ADDR"), 50, 10, "3m")
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer db.Close()
	fmt.Printf("[INFO] Successfully connected to database\n")

	store := store.NewStorage(db, log)

	server := &httpServer{
		db:     db,
		logger: log,
		store:  store,
	}

	// cron
	server.startCronUpdSubscriptions()
	server.startCronUpdUserSubscriptionBeforEndSub()
	server.startCronAutoActivateExpiredTransactions()

	// start http server
	if err := server.httpStart(); err != nil {
		log.Error(err.Error())
	}
}
