package main

import (
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/stdi0/worker_binance_announce/config"
	"github.com/stdi0/worker_binance_announce/database"
	"github.com/stdi0/worker_binance_announce/worker"
)

func main() {
	log.Println("Getting database config...")
	dbConf := config.GetDatabaseConfig()

	log.Println("Getting worker config...")
	workerConf := config.GetWorkerConfig()

	log.Println("Getting trader config...")
	traderConf := config.GetTraderConfig()

	var db *database.Database

	for {
		var err error
		log.Println("Trying connect to database...")
		db, err = database.NewDatabase(dbConf)
		if err != nil {
			log.Printf("[ERROR] %s", err.Error())
			log.Println("Trying reconnect after 3 sec...")
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}

	fetcher := worker.NewFetcher(db)
	trader := worker.NewTrader(traderConf)
	w := worker.NewWorker(workerConf, fetcher, trader)
	w.Start()
}
