package repository

import (
	"Text-Gathering-Service/misc"
	"Text-Gathering-Service/models"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func New() *Database {
	file := misc.Must(sql.Open("sqlite3", "./internal/repository/database.db"))
	misc.Must(file.Exec(
		`CREATE TABLE IF NOT EXISTS NewGatheredText(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			text TEXT NOT NULL,
			audioText TEXT NOT NULL,
			isCheck INTEGER NOT NULL DEFAULT 0
		)`))
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
	query := `SELECT EXISTS(SELECT 1 FROM NewGatheredText WHERE text = ?)`
	d.db.QueryRow(query, t).Scan(&exists)

	return exists
}

func (d *Database) StoreIntoDB(text string, status bool) bool {
	if !d._checkExistDatas(text) {
		misc.Must(d.db.Exec(`INSERT INTO NewGatheredText(text, audioText, isCheck) 
							VALUES (?, ?, ?)`,
			text, "internal/repository/wait_clips/"+text+".wav", status))
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

func (d *Database) GetAllWaitClipsDatas() []models.ResCheckedDatas {
	rows := misc.Must(d.db.Query(`
		SELECT id, text, audioText 
		FROM NewGatheredText 
		WHERE isCheck = 0`))
	defer rows.Close()

	var datas []models.ResCheckedDatas
	var rem_id []int64

	for rows.Next() {
		var data models.ResCheckedDatas

		err := rows.Scan(&data.ID, &data.Text, &data.Voice)
		if err != nil {
			panic(err)
		}

		if _, err := os.Open(data.Voice); err != nil && os.IsNotExist(err) {
			_ = os.Remove(data.Voice)
			rem_id = append(rem_id, data.ID)
			fmt.Printf("ID:%d Text:%s Path:%s\n", data.ID, data.Text, data.Voice)
			log.Printf("This path '%s' was deleted", data.Voice)
			continue
		}
		datas = append(datas, data)
	}

	for _, id := range rem_id {
		misc.Must(d.db.Exec(`DELETE FROM NewGatheredText WHERE id = ?`, id))
	}

	return datas
}

func (d *Database) ChangeStatusClipDatas(id int64, text string, status bool) {
	var path string
	if status {
		path = "internal/repository/successful_clips/"
		os.MkdirAll(path, os.ModePerm)
		loc_path := path + text + ".wav"
		os.Rename("internal/repository/wait_clips/"+text+".wav", loc_path)
		misc.Must(d.db.Exec(`
			UPDATE NewGatheredText 
			SET isCheck = 1,
				text = ?,
				audioText = ?
			WHERE id = ?`,
			text, loc_path, id))
	} else {
		path = "internal/repository/not_successful_clips/"
		os.MkdirAll(path, os.ModePerm)
		loc_path := path + text + ".wav"
		os.Rename("internal/repository/wait_clips/"+text+".wav", loc_path)
		misc.Must(d.db.Exec(`
			UPDATE NewGatheredText 
			SET isCheck = -1,
				audioText = ?
			WHERE id = ?`,
			loc_path, id))
	}
	defer d.db.Close()
}
