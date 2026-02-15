package load_test

import (
	"database/sql"
	"log"

	"github.com/chrwhy/simple/examples/go/util"
)

func CreateExternalTable(db *sql.DB) error {
	createSQL := `CREATE TABLE t1 (a INTEGER PRIMARY KEY, b, c, d);`
	createFtsTableSQL := `CREATE VIRTUAL TABLE ft USING fts5(b, c, d UNINDEXED, content=t1, content_rowid=a, tokenize = 'simple 1');`
	triggerSQL := `CREATE TRIGGER t1_ai AFTER INSERT ON t1 BEGIN
  					INSERT INTO ft(rowid, b, c,d) VALUES (new.a, new.b, new.c, new.d);
					END;`
	idxSQL := `CREATE INDEX idx_fts_d ON t1(d)`

	_, err := db.Exec(createSQL)
	_, err = db.Exec(createFtsTableSQL)
	_, err = db.Exec(triggerSQL)
	_, err = db.Exec(idxSQL)
	if err != nil {
		log.Printf("CreateChatGroupTable error: %v", err)
	}
	return err
}

func MockData(db *sql.DB) error {
	for i := 0; i < 9999999; i++ {
		a := i
		b := util.RandomFirstName()
		c := util.RandomLastName()
		d := util.RandInt(0, 9999999999)
		insertSQL := `INSERT INTO t1(a,b,c,d) VALUES (?, ?, ?, ?);`
		_, err := db.Exec(insertSQL, a, b, c, d)
		if err != nil {
			log.Printf("InsertChatGroup error: %v", err)
		}
	}
	return nil
}
