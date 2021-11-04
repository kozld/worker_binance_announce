package database

const CreateTableQuery = `CREATE TABLE IF NOT EXISTS binance_api (
	hash bytea UNIQUE NOT NULL,
	text text NOT NULL,
	time timestamp with time zone default current_timestamp,
	PRIMARY KEY (hash)
);`

const SelectQuery = `SELECT text FROM binance_api ORDER BY time DESC LIMIT 10;`
