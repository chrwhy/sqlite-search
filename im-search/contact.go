package im_search

import (
	"database/sql"
	"log"
)

type Contact struct {
	Uid   int
	Name  string
	Alias string
}

// CreateContactTable creates the FTS5 virtual table if it doesn't exist.
func CreateContactTable(db *sql.DB) error {
	createSQL := `CREATE VIRTUAL TABLE IF NOT EXISTS contact USING fts5(uid, name, alias, tokenize = 'simple 1');`
	_, err := db.Exec(createSQL)
	if err != nil {
		log.Printf("CreateContactTable error: %v", err)
	}
	return err
}

// InsertContact inserts a new contact record.
func InsertContact(db *sql.DB, c Contact) error {
	insertSQL := `INSERT INTO contact(uid, name, alias) VALUES (?, ?, ?);`
	_, err := db.Exec(insertSQL, c.Uid, c.Name, c.Alias)
	if err != nil {
		log.Printf("InsertContact error: %v", err)
	}
	return err
}

// UpdateContact updates name and alias for an existing uid.
func UpdateContact(db *sql.DB, c Contact) error {
	updateSQL := `UPDATE contact SET name = ?, alias = ? WHERE uid = ?;`
	res, err := db.Exec(updateSQL, c.Name, c.Alias, c.Uid)
	if err != nil {
		log.Printf("UpdateContact error: %v", err)
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		log.Printf("UpdateContact: no rows updated for uid=%d", c.Uid)
	}
	return nil
}

// DeleteContact removes a contact by uid.
func DeleteContact(db *sql.DB, uid int) error {
	deleteSQL := `DELETE FROM contact WHERE uid = ?;`
	_, err := db.Exec(deleteSQL, uid)
	if err != nil {
		log.Printf("DeleteContact error: %v", err)
	}
	return err
}

// GetContact retrieves a single contact by uid.
func GetContact(db *sql.DB, uid int) (Contact, error) {
	var c Contact
	query := `SELECT uid, name, alias FROM contact WHERE uid = ? LIMIT 1;`
	row := db.QueryRow(query, uid)
	err := row.Scan(&c.Uid, &c.Name, &c.Alias)
	if err == sql.ErrNoRows {
		return c, nil
	}
	if err != nil {
		log.Printf("GetContact error: %v", err)
	}
	return c, err
}

// SearchContacts uses FTS5 MATCH to find matching contacts.
// Provide a raw FTS5 query like: "name:alice OR alias:bob" or simple term "alice".
func SearchContacts(db *sql.DB, clause string) ([]Contact, error) {
	sqlStmt := "SELECT uid, simple_highlight(contact, 1, '[', ']') , simple_highlight(contact, 2, '[', ']') FROM contact WHERE (name MATCH ('" + clause + "')) OR (alias MATCH ('" + clause + "'));"
	rows, err := db.Query(sqlStmt)
	if err != nil {
		log.Printf("SearchContacts query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var results []Contact
	for rows.Next() {
		var c Contact
		if err := rows.Scan(&c.Uid, &c.Name, &c.Alias); err != nil {
			log.Printf("SearchContacts scan error: %v", err)
			continue
		}
		results = append(results, c)
	}
	if err := rows.Err(); err != nil {
		log.Printf("SearchContacts rows error: %v", err)
		return results, err
	}
	return results, nil
}

// SeedContacts inserts example contacts (friends) with Chinese names and pinyin aliases.
// Safe to call multiple times; individual insert errors are logged but not fatal.
func SeedContacts(db *sql.DB) error {
	contacts := []Contact{
		{Uid: 1001, Name: "张三", Alias: "zhangsan"},
		{Uid: 1002, Name: "李四", Alias: "lisi"},
		{Uid: 1003, Name: "王五", Alias: "wangwu"},
		{Uid: 1004, Name: "赵六", Alias: "zhaoliu"},
		{Uid: 1005, Name: "孙晓明", Alias: "sunxiaoming"},
		{Uid: 1006, Name: "周杰伦", Alias: "zhoujielun"},
		{Uid: 1007, Name: "陈奕迅", Alias: "chenyixun"},
		{Uid: 1008, Name: "小红", Alias: "xiaohong"},
	}

	for _, c := range contacts {
		if err := InsertContact(db, c); err != nil {
			log.Printf("SeedContacts: failed to insert uid=%d name=%q: %v", c.Uid, c.Name, err)
		}
	}
	return nil
}
