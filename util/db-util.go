package util

import (
	"database/sql"
	"fmt"
	"log"
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

func Query(db *sql.DB, querySQL string) {
	t0 := time.Now()
	rows, err := db.Query(querySQL)
	log.Println("Query cost: ", time.Since(t0))
	if err != nil {
		log.Printf("query error: %v with sql: %s", err, querySQL)
		return
	}
	defer rows.Close()
	t1 := time.Now()

	log.Println("Timestamp 0: ", time.Now().UnixMilli())
	cols, err := rows.Columns()
	if err != nil {
		log.Printf("columns error: %v", err)
		return
	}
	colCount := len(cols)

	log.Println("Timestamp 1: ", time.Now().UnixMilli())
	// prepare a slice of *sql.NullString to accept nullable column values
	values := make([]interface{}, colCount)
	for i := range values {
		values[i] = new(sql.NullString)
	}

	log.Println("Timestamp 2: ", time.Now().UnixMilli())
	log.Println(rows)
	rowIndex := 0
	for rows.Next() {
		log.Println("Timestamp 2-000: ", time.Now().UnixMilli())
		if err := rows.Scan(values...); err != nil {
			log.Println("Timestamp 2-0: ", time.Now().UnixMilli())
			log.Printf("scan error: %v", err)
			continue
		}

		log.Println("Timestamp 2-1: ", time.Now().UnixMilli())
		pairs := make([]string, colCount)
		for i := 0; i < colCount; i++ {
			ns := values[i].(*sql.NullString)
			val := "NULL"
			if ns.Valid {
				val = ns.String
			}
			pairs[i] = fmt.Sprintf("%s=%s", cols[i], val)
		}
		log.Println("Timestamp 2-2: ", time.Now().UnixMilli())
		log.Printf("Row %d: %s", rowIndex, fmt.Sprintf("%v", pairs))
		rowIndex++
	}

	log.Println("Timestamp 3: ", time.Now().UnixMilli())
	if err := rows.Err(); err != nil {
		log.Printf("rows iteration error: %v", err)
	}
	log.Println("Timestamp 4: ", time.Now().UnixMilli())
	log.Println("Print cost: ", time.Since(t1))
	log.Printf("Total : %d\n", rowIndex)
}
