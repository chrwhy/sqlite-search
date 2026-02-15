package im_search

import (
	"database/sql"
	"log"
)

type GroupMember struct {
	Gid          int
	Uid          int
	Name         string
	Alias        string
	AliasInGroup string
}

// CreateGroupMemberTable creates the FTS5 virtual table if it doesn't exist.
func CreateGroupMemberTable(db *sql.DB) error {
	createSQL := `CREATE VIRTUAL TABLE IF NOT EXISTS group_member USING fts5(gid, uid, name, alias, alias_in_group, tokenize = 'simple 1');`
	_, err := db.Exec(createSQL)
	if err != nil {
		log.Printf("CreateGroupMemberTable error: %v", err)
	}
	return err
}

// InsertGroupMember inserts a new group member record.
func InsertGroupMember(db *sql.DB, gm GroupMember) error {
	insertSQL := `INSERT INTO group_member(gid, uid, name, alias, alias_in_group) VALUES (?, ?, ?, ?, ?);`
	_, err := db.Exec(insertSQL, gm.Gid, gm.Uid, gm.Name, gm.Alias, gm.AliasInGroup)
	if err != nil {
		log.Printf("InsertGroupMember error: %v", err)
	}
	return err
}

// UpdateGroupMember updates name, alias and alias_in_group for an existing gid+uid.
func UpdateGroupMember(db *sql.DB, gm GroupMember) error {
	updateSQL := `UPDATE group_member SET name = ?, alias = ?, alias_in_group = ? WHERE gid = ? AND uid = ?;`
	res, err := db.Exec(updateSQL, gm.Name, gm.Alias, gm.AliasInGroup, gm.Gid, gm.Uid)
	if err != nil {
		log.Printf("UpdateGroupMember error: %v", err)
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		log.Printf("UpdateGroupMember: no rows updated for gid=%d uid=%d", gm.Gid, gm.Uid)
	}
	return nil
}

// DeleteGroupMember removes a group member by gid and uid.
func DeleteGroupMember(db *sql.DB, gid, uid int) error {
	deleteSQL := `DELETE FROM group_member WHERE gid = ? AND uid = ?;`
	_, err := db.Exec(deleteSQL, gid, uid)
	if err != nil {
		log.Printf("DeleteGroupMember error: %v", err)
	}
	return err
}

// GetGroupMember retrieves a single group member by gid and uid.
func GetGroupMember(db *sql.DB, gid, uid int) (GroupMember, error) {
	var gm GroupMember
	query := `SELECT gid, uid, name, alias, alias_in_group FROM group_member WHERE gid = ? AND uid = ? LIMIT 1;`
	row := db.QueryRow(query, gid, uid)
	err := row.Scan(&gm.Gid, &gm.Uid, &gm.Name, &gm.Alias, &gm.AliasInGroup)
	if err == sql.ErrNoRows {
		return gm, nil
	}
	if err != nil {
		log.Printf("GetGroupMember error: %v", err)
	}
	return gm, err
}

// SearchGroupMembers uses FTS5 MATCH to find matching group members.
// Provide a raw FTS5 query like: "name:alice OR alias:bob" or a simple term "alice".
func SearchGroupMembers(db *sql.DB, clause string) ([]GroupMember, error) {
	sqlStmt := "SELECT gid, uid, simple_highlight(group_member, 2, '[', ']') , simple_highlight(group_member, 3, '[', ']'), simple_highlight(group_member, 4, '[', ']') FROM group_member WHERE (name MATCH ('" + clause + "')) OR (alias MATCH ('" + clause + "')) OR (alias_in_group MATCH ('" + clause + "'));"
	log.Println(sqlStmt)
	rows, err := db.Query(sqlStmt)
	if err != nil {
		log.Printf("SearchGroupMembers query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var results []GroupMember
	for rows.Next() {
		var gm GroupMember
		if err := rows.Scan(&gm.Gid, &gm.Uid, &gm.Name, &gm.Alias, &gm.AliasInGroup); err != nil {
			log.Printf("SearchGroupMembers scan error: %v", err)
			continue
		}
		results = append(results, gm)
	}
	if err := rows.Err(); err != nil {
		log.Printf("SearchGroupMembers rows error: %v", err)
		return results, err
	}
	return results, nil
}

// SeedGroupMembers inserts example group members for seeded chat groups.
// Safe to call multiple times; individual insert errors are logged but not fatal.
func SeedGroupMembers(db *sql.DB) error {
	members := []GroupMember{
		// Members for group 1 (开发组)
		{Gid: 1, Uid: 101, Name: "李雷", Alias: "lilei", AliasInGroup: "李"},
		{Gid: 1, Uid: 102, Name: "王强", Alias: "wangqiang", AliasInGroup: "王"},
		// Members for group 2 (产品讨论)
		{Gid: 2, Uid: 201, Name: "张伟", Alias: "zhangwei", AliasInGroup: "张"},
		{Gid: 2, Uid: 202, Name: "陈静", Alias: "chenjing", AliasInGroup: "陈"},
		// Members for group 3 (市场推广)
		{Gid: 3, Uid: 301, Name: "刘洋", Alias: "liuyang", AliasInGroup: "刘"},
		// Members for group 4 (运维组)
		{Gid: 4, Uid: 401, Name: "赵磊", Alias: "zhaolei", AliasInGroup: "赵"},
		// Members for group 5 (设计团队)
		{Gid: 5, Uid: 501, Name: "韩梅梅", Alias: "hanmeimei", AliasInGroup: "韩"},
		// Members for group 1001 (Friends)
		{Gid: 1001, Uid: 10001, Name: "小红", Alias: "xiaohong", AliasInGroup: "小红"},
		{Gid: 1001, Uid: 10002, Name: "小明", Alias: "xiaoming", AliasInGroup: "小明2"},
	}

	for _, m := range members {
		if err := InsertGroupMember(db, m); err != nil {
			log.Printf("SeedGroupMembers: failed to insert gid=%d uid=%d name=%q: %v", m.Gid, m.Uid, m.Name, err)
		}
	}
	return nil
}
