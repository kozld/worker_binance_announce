package worker

import (
	"log"
	"time"

	"github.com/stdi0/worker_binance_announce/config"
)

type Worker struct {
	conf           *config.WorkerConfig
	fetcher        *Fetcher
	excludedTokens map[string]bool
}

func NewWorker(conf *config.WorkerConfig, fetcher *Fetcher) *Worker {
	return &Worker{conf, fetcher, make(map[string]bool)}
}

func (w *Worker) Start() {
	for {
		log.Println("Fetching announcements...")

		tokens := w.fetcher.Fetch()
		for _, gem := range tokens {
			if _, exist := w.excludedTokens[gem]; !exist {
				log.Printf("[TOKEN] (%s)", gem)
				w.excludedTokens[gem] = true
			} else {
				log.Printf("Token (%s) already processed", gem)
			}
		}

		log.Println("Sleep 3 sec...")
		time.Sleep(3 * time.Second)
	}
}
