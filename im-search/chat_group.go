package im_search

import (
	"database/sql"
	"log"
)

type ChatGroup struct {
	Gid   int
	Name  string
	Alias string
}

// CreateChatGroupTable creates the FTS5 virtual table if it doesn't exist.
func CreateChatGroupTable(db *sql.DB) error {
	createSQL := `CREATE VIRTUAL TABLE IF NOT EXISTS chat_group USING fts5(gid, name, alias, tokenize = 'simple 1');`
	_, err := db.Exec(createSQL)
	if err != nil {
		log.Printf("CreateChatGroupTable error: %v", err)
	}
	return err
}

// InsertChatGroup inserts a new chat group record.
func InsertChatGroup(db *sql.DB, g ChatGroup) error {
	insertSQL := `INSERT INTO chat_group(gid, name, alias) VALUES (?, ?, ?);`
	_, err := db.Exec(insertSQL, g.Gid, g.Name, g.Alias)
	if err != nil {
		log.Printf("InsertChatGroup error: %v", err)
	}
	return err
}

// UpdateChatGroup updates name and alias for an existing gid.
func UpdateChatGroup(db *sql.DB, g ChatGroup) error {
	updateSQL := `UPDATE chat_group SET name = ?, alias = ? WHERE gid = ?;`
	res, err := db.Exec(updateSQL, g.Name, g.Alias, g.Gid)
	if err != nil {
		log.Printf("UpdateChatGroup error: %v", err)
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		log.Printf("UpdateChatGroup: no rows updated for gid=%d", g.Gid)
	}
	return nil
}

// DeleteChatGroup removes a chat group by gid.
func DeleteChatGroup(db *sql.DB, gid int) error {
	deleteSQL := `DELETE FROM chat_group WHERE gid = ?;`
	_, err := db.Exec(deleteSQL, gid)
	if err != nil {
		log.Printf("DeleteChatGroup error: %v", err)
	}
	return err
}

// GetChatGroup retrieves a single chat group by gid.
func GetChatGroup(db *sql.DB, gid int) (ChatGroup, error) {
	var g ChatGroup
	query := `SELECT gid, name, alias FROM chat_group WHERE gid = ? LIMIT 1;`
	row := db.QueryRow(query, gid)
	err := row.Scan(&g.Gid, &g.Name, &g.Alias)
	if err == sql.ErrNoRows {
		return g, nil
	}
	if err != nil {
		log.Printf("GetChatGroup error: %v", err)
	}
	return g, err
}

// SearchChatGroups uses FTS5 MATCH to find matching chat groups.
// Provide a raw FTS5 query like: "name:alice OR alias:bob" or simple term "alice".
func SearchChatGroups(db *sql.DB, clause string) ([]ChatGroup, error) {
	sqlStmt := "SELECT gid, simple_highlight(chat_group, 1, '[', ']') , simple_highlight(chat_group, 2, '[', ']') FROM chat_group WHERE (name MATCH ('" + clause + "')) OR (alias MATCH ('" + clause + "'));"
	log.Println(sqlStmt)
	rows, err := db.Query(sqlStmt)
	if err != nil {
		log.Printf("SearchChatGroups query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var results []ChatGroup
	for rows.Next() {
		var g ChatGroup
		if err := rows.Scan(&g.Gid, &g.Name, &g.Alias); err != nil {
			log.Printf("SearchChatGroups scan error: %v", err)
			continue
		}
		results = append(results, g)
	}
	if err := rows.Err(); err != nil {
		log.Printf("SearchChatGroups rows error: %v", err)
		return results, err
	}
	return results, nil
}

// SeedChatGroups inserts a small set of initial chat groups for examples and testing.
// It is safe to call multiple times; errors on individual inserts are logged but not fatal.
func SeedChatGroups(db *sql.DB) error {
	groups := []ChatGroup{
		{Gid: 1, Name: "开发组", Alias: "dev"},
		{Gid: 2, Name: "产品讨论", Alias: "product"},
		{Gid: 3, Name: "市场推广", Alias: "marketing"},
		{Gid: 4, Name: "运维组", Alias: "ops"},
		{Gid: 5, Name: "设计团队", Alias: "design"},
		{Gid: 1001, Name: "Friends", Alias: "friends"},
	}

	for _, g := range groups {
		if err := InsertChatGroup(db, g); err != nil {
			// Log and continue; seeding should not fail the whole app if one insert errors.
			log.Printf("SeedChatGroups: failed to insert gid=%d name=%q alias=%q: %v", g.Gid, g.Name, g.Alias, err)
		}
	}
	return nil
}
