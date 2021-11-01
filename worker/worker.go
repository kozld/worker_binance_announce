package worker

import (
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

		if len(tokens) == 0 {
			log.Println("No announcements yet")
		}

		for _, gem := range tokens {
			if _, exist := w.excludedTokens[gem]; !exist {
				log.Printf("!!! NEW TOKEN (%s) FOUND !!!", gem)

				log.Println("Trying to buy on Gate.io...")
				log.Println(w.trader.CreateOrder(gem, false))

				w.excludedTokens[gem] = true
			} else {
				log.Printf("Token (%s) already processed", gem)
			}
		}

		log.Println("Sleep 3 sec...")
		time.Sleep(3 * time.Second)
	}
}
