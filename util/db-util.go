package util

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand/v2"
	"time"

	"github.com/mattn/go-sqlite3"
)

func InitDB() *sql.DB {
	sql.Register("sqlite3_simple",
		&sqlite3.SQLiteDriver{
			Extensions: []string{
				"./libsimple-osx-x64/libsimple",
			},
		})

	//db, err := sql.Open("sqlite3_simple", ":memory:")
	db, err := sql.Open("sqlite3_simple", "example.db")
	if err != nil {
		log.Fatalf("open error: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("ping error: ", err)
	}
	return db
}

func CreateTable(db *sql.DB) {
	createTableSQL := `CREATE VIRTUAL TABLE t1 USING fts5(biz_id, text, tokenize = 'simple 1');`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Table created successfully")
}

func InsertRecord(db *sql.DB, bizId int, text string) {
	if bizId <= 0 {
		bizId = rand.Int()
	}
	insertSQL := `INSERT INTO t1(biz_id, text) VALUES (?, ?)`
	_, err := db.Exec(insertSQL, bizId, text)
	if err != nil {
		log.Fatal(err)
	}
	//log.Println("Records inserted successfully")
}

func DeleteRecord(db *sql.DB, bizId int) {
	if bizId <= 0 {
		return
	}
	insertSQL := `DELETE FROM t1 WHERE biz_id = ?`
	_, err := db.Exec(insertSQL, bizId)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Record deleted successfully")
}

func Query(db *sql.DB, querySQL string) {
	t0 := time.Now()
	rows, err := db.Query(querySQL)
	log.Println("Query cost: ", time.Since(t0))
	if err != nil {
		log.Printf("query error: %v with sql: %s", err, querySQL)
		return
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		log.Printf("columns error: %v", err)
		return
	}
	colCount := len(cols)

	// prepare a slice of *sql.NullString to accept nullable column values
	values := make([]interface{}, colCount)
	for i := range values {
		values[i] = new(sql.NullString)
	}

	rowIndex := 0
	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			log.Printf("scan error: %v", err)
			continue
		}

		pairs := make([]string, colCount)
		for i := 0; i < colCount; i++ {
			ns := values[i].(*sql.NullString)
			val := "NULL"
			if ns.Valid {
				val = ns.String
			}
			pairs[i] = fmt.Sprintf("%s=%s", cols[i], val)
		}
		log.Printf("Row %d: %s", rowIndex, fmt.Sprintf("%v", pairs))
		rowIndex++
	}

	if err := rows.Err(); err != nil {
		log.Printf("rows iteration error: %v", err)
	}
	log.Printf("Total : %d\n", rowIndex)
}
