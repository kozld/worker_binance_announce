package worker

import (
	"log"
	"regexp"

	"github.com/stdi0/worker_binance_announce/database"
)

const RegularExpression = `Will List.*?\(([[:upper:]]*)\)`

type Fetcher struct {
	db *database.Database
}

func NewFetcher(db *database.Database) *Fetcher {
	// Create database table if not exist
	log.Println("Creating postgres table if not exist...")
	db.Conn.Exec(database.CreateTableQuery)

	return &Fetcher{db}
}

func (f *Fetcher) Fetch() []string {
	tokens := make([]string, 0)
	re := regexp.MustCompile(RegularExpression)
	rows, _ := f.db.Conn.Query(database.SelectQuery)

	for {
		var text string
		rows.Scan(&text)

		match := re.FindStringSubmatch(text)
		if len(match) > 1 {
			tokens = append(tokens, match[1])
		}

		if !rows.Next() {
			break
		}
	}

	return tokens
}
