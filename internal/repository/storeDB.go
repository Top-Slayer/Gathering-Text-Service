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
	query := `SELECT EXISTS (
				SELECT 1 
				FROM NewGatheredText 
				WHERE text = ?
			)`
	d.db.QueryRow(query, t).Scan(&exists)

	return exists
}

func (d *Database) StoreIntoDB(text string, status bool) bool {
	if !d._checkExistDatas(text) {
		latest_id := d.GetLatestID() + 1
		path := fmt.Sprintf("internal/repository/wait_clips/voice_id_%d.wav", latest_id)
		misc.Must(d.db.Exec(`INSERT INTO NewGatheredText
							(id, text, audioText, isCheck) 
							VALUES (?, ?, ?, ?)`,
			latest_id, text, path, status))

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
	var old_text string

	err := d.db.QueryRow("SELECT text FROM NewGatheredText WHERE id = ?", id).Scan(&old_text)
	if err != nil {
		log.Println(err)
	}

	if status {
		path = "internal/repository/successful_clips/"
	} else {
		path = "internal/repository/not_successful_clips/"
	}

	os.MkdirAll(path, os.ModePerm)
	id_voice := fmt.Sprintf("voice_id_%d.wav", id)
	loc_path := fmt.Sprintf("internal/repository/wait_clips/%s", id_voice)
	des_path := fmt.Sprintf("%s%s", path, id_voice)

	var isCheck int8

	if status {
		isCheck = 1
	} else {
		isCheck = -1
	}

	os.Rename(loc_path, des_path)
	misc.Must(d.db.Exec(`
			UPDATE NewGatheredText 
			SET isCheck = ?,
				text = ?,
				audioText = ?
			WHERE id = ?`,
		isCheck, text, des_path, id))
}

func (d *Database) GetLatestID() int64 {
	var id int64
	err := d.db.QueryRow("SELECT MAX(id) FROM NewGatheredText", id).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0
	}
	return id
}

func (d *Database) CloseDatabase() {
	d.db.Close()
}
