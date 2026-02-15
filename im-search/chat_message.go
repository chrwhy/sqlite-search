package im_search

import (
	"database/sql"
	"log"
	"strings"
)

type ChatMessage struct {
	Cid         int
	SubjectId   int
	SubjectType string
	Message     string
}

// CreateChatMessageTable creates the FTS5 virtual table if it doesn't exist.
func CreateChatMessageTable(db *sql.DB) error {
	createSQL := `CREATE VIRTUAL TABLE IF NOT EXISTS chat_message USING fts5(cid, subject_id, subject_type, message, tokenize = 'simple 1');`
	_, err := db.Exec(createSQL)
	if err != nil {
		log.Printf("CreateChatMessageTable error: %v", err)
	}
	return err
}

// InsertChatMessage inserts a new chat message record.
func InsertChatMessage(db *sql.DB, m ChatMessage) error {
	insertSQL := `INSERT INTO chat_message(cid, subject_id, subject_type, message) VALUES (?, ?, ?, ?);`
	_, err := db.Exec(insertSQL, m.Cid, m.SubjectId, m.SubjectType, m.Message)
	if err != nil {
		log.Printf("InsertChatMessage error: %v", err)
	}
	return err
}

// UpdateChatMessage updates subject and message fields for an existing cid.
func UpdateChatMessage(db *sql.DB, m ChatMessage) error {
	updateSQL := `UPDATE chat_message SET subject_id = ?, subject_type = ?, message = ? WHERE cid = ?;`
	res, err := db.Exec(updateSQL, m.SubjectId, m.SubjectType, m.Message, m.Cid)
	if err != nil {
		log.Printf("UpdateChatMessage error: %v", err)
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		log.Printf("UpdateChatMessage: no rows updated for cid=%d", m.Cid)
	}
	return nil
}

// DeleteChatMessage removes a chat message by cid.
func DeleteChatMessage(db *sql.DB, cid int) error {
	deleteSQL := `DELETE FROM chat_message WHERE cid = ?;`
	_, err := db.Exec(deleteSQL, cid)
	if err != nil {
		log.Printf("DeleteChatMessage error: %v", err)
	}
	return err
}

// GetChatMessage retrieves a single chat message by cid.
func GetChatMessage(db *sql.DB, cid int) (ChatMessage, error) {
	var m ChatMessage
	query := `SELECT cid, subject_id, subject_type, message FROM chat_message WHERE cid = ? LIMIT 1;`
	row := db.QueryRow(query, cid)
	err := row.Scan(&m.Cid, &m.SubjectId, &m.SubjectType, &m.Message)
	if err == sql.ErrNoRows {
		return m, nil
	}
	if err != nil {
		log.Printf("GetChatMessage error: %v", err)
	}
	return m, err
}

func SearchChatMessages(db *sql.DB, q string) ([]ChatMessage, error) {
	terms := strings.Fields(q)
	if len(terms) == 0 {
		return nil, nil
	}

	for i, t := range terms {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		// Escape double quotes inside term and wrap in quotes for phrase search
		t = strings.ReplaceAll(t, `"`, `""`)
		terms[i] = `"` + t + `"`
	}

	matchExpr := strings.Join(terms, " AND ")
	sqlStmt := `SELECT cid, subject_id, subject_type, simple_highlight(chat_message, 3, '[', ']') FROM chat_message WHERE message MATCH ?;`
	rows, err := db.Query(sqlStmt, matchExpr)
	if err != nil {
		log.Printf("SearchChatMessages query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var results []ChatMessage
	for rows.Next() {
		var m ChatMessage
		if err := rows.Scan(&m.Cid, &m.SubjectId, &m.SubjectType, &m.Message); err != nil {
			log.Printf("SearchChatMessages scan error: %v", err)
			continue
		}
		results = append(results, m)
	}
	if err := rows.Err(); err != nil {
		log.Printf("SearchChatMessages rows error: %v", err)
		return results, err
	}
	return results, nil
}

// SeedChatMessages inserts example chat messages for both friend chats (subject_type="user")
// and group chats (subject_type="group"). It is safe to call multiple times; individual
// insert errors are logged but not fatal.
func SeedChatMessages(db *sql.DB) error {
	messages := []ChatMessage{
		// Direct friend messages (subject_type = "contact"), subject_id is friend's uid
		{Cid: 10001, SubjectId: 1001, SubjectType: "contact", Message: "你好，张三！最近怎么样？"},
		{Cid: 10002, SubjectId: 1002, SubjectType: "contact", Message: "李四，明天一起吃饭吗？"},
		{Cid: 10003, SubjectId: 1006, SubjectType: "contact", Message: "周杰伦的歌真好听。"},

		// Group messages (subject_type = "group"), subject_id is gid
		{Cid: 20001, SubjectId: 1, SubjectType: "group", Message: "大家好，今天的代码 review 安排在下午 3 点。"},
		{Cid: 20002, SubjectId: 1, SubjectType: "group", Message: "请把你负责的模块 checklist 发一下。"},
		{Cid: 20003, SubjectId: 2, SubjectType: "group", Message: "产品需求已经更新，请查看文档。"},
		{Cid: 20004, SubjectId: 1001, SubjectType: "group", Message: "今晚聚餐地点：老地方。谁能来请回复。"},
	}

	for _, m := range messages {
		if err := InsertChatMessage(db, m); err != nil {
			log.Printf("SeedChatMessages: failed to insert cid=%d subject=%d type=%s: %v", m.Cid, m.SubjectId, m.SubjectType, err)
		}
	}
	return nil
}
