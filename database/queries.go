package database

const CreateTableQuery = `CREATE TABLE IF NOT EXISTS binance_page (
	hash bytea UNIQUE NOT NULL,
	text text NOT NULL,
	time timestamp with time zone default current_timestamp,
	PRIMARY KEY (hash)
);`

const SelectQuery = `SELECT text FROM binance_page ORDER BY time ASC LIMIT 10;`
