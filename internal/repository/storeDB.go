package repository

import (
	"Text-Gathering-Service/misc"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func New() *Database {
	file := misc.Must(sql.Open("sqlite3", "./internal/repository/database.db"))
	misc.Must(file.Exec(`CREATE TABLE IF NOT EXISTS GatheredText(text TEXT NOT NULL)`))

	return &Database{
		db: file,
	}
}

// check that word or sentense already have or not

func (d *Database) StoreIntoDB(text string) {
	misc.Must(d.db.Exec("INSERT INTO GatheredText(text) VALUES (?)", text))
}

func (d *Database) Close() {
	defer d.db.Close()
}
