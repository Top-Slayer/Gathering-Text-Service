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
	misc.Must(file.Exec(
		`CREATE TABLE IF NOT EXISTS Categories(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(30) NOT NULL
		)`))

	return &Database{
		db: file,
	}
}

func (d *Database) _checkExistDatas(t string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM GatheredText WHERE text = ?)`
	d.db.QueryRow(query, t).Scan(&exists)

	return exists
}

func (d *Database) StoreIntoDB(text string) bool {
	if !d._checkExistDatas(text) {
		misc.Must(d.db.Exec(`INSERT INTO GatheredText(text) VALUES (?)`, text))
		defer d.db.Close()

		return true
	} else {

		return false
	}
}

func (d *Database) GetAllCategoryDatas() []map[string]interface{} {
	rows := misc.Must(d.db.Query(`SELECT name FROM Categories`))
	defer rows.Close()

	var res []map[string]interface{}
	columns := misc.Must(rows.Columns())

	for rows.Next() {
		columnPointers := make([]interface{}, len(columns))
		columnValues := make([]interface{}, len(columns))

		for i := range columnPointers {
			columnPointers[i] = &columnValues[i]
		}

		err := rows.Scan(columnPointers...)
		if err != nil {
			return nil
		}

		rowMap := make(map[string]interface{})
		for i, colName := range columns {
			rowMap[colName] = columnValues[i]
		}
		res = append(res, rowMap)
	}

	return res
}
