package worker

import (
	"fmt"
	"log"
	"time"

	"github.com/stdi0/worker_binance_announce/config"
)

type Worker struct {
	conf           *config.WorkerConfig
	fetcher        *Fetcher
	trader         *Trader
	excludedTokens map[string]bool
}

func NewWorker(conf *config.WorkerConfig, fetcher *Fetcher, trader *Trader) *Worker {
	// Exclude "CHESS" token
	excludeTokens := make(map[string]bool)
	excludeTokens["CHESS"] = true
	return &Worker{conf, fetcher, trader, excludeTokens}
}

func (w *Worker) Start() {
	for {
		log.Println("Fetching announcements...")

		tokens := w.fetcher.Fetch()
		for _, gem := range tokens {
			if _, exist := w.excludedTokens[gem]; !exist {
				log.Printf("[TOKEN] (%s)", gem)
				fmt.Println("[TRADE]", w.trader.CreateOrder(gem, false))
				w.excludedTokens[gem] = true
			} else {
				log.Printf("Token (%s) already processed", gem)
			}
		}

		log.Println("Sleep 3 sec...")
		time.Sleep(3 * time.Second)
	}
}
