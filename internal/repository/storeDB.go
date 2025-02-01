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

func (d *Database) _checkExistDatas(t string) bool {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM GatheredText WHERE text = ?)"
	d.db.QueryRow(query, t).Scan(&exists)

	return exists
}

func (d *Database) StoreIntoDB(text string) bool {
	if !d._checkExistDatas(text) {
		misc.Must(d.db.Exec("INSERT INTO GatheredText(text) VALUES (?)", text))
		defer d.db.Close()

		return true
	} else {

		return false
	}
}
