package spotlight

import (
	"database/sql"
	"log"
	"math/rand"
	"os"
)

func InitData(db *sql.DB) {
	records := []string{
		"周杰伦 Jay Chou: \"最美的不是下雨天，是曾与你躲过雨的屋檐\"",
		"I love China! 我爱中国!",
		"@English &special _characters.\"''bacon-&and''-eggs%",
		"中转箱",
		"李安",
		"西安",
		"周星驰",
		"李会",
		"刘慧子",
		"刘亚男",
		"张蔷",
		"张倩歌",
		"张强哥",
		"练习",
		"吕布",
		"绿色",
		"驴子",
		"13825638962",
		"珠海@中国",
		"北京@中國",
		"living",
		"中華人民共和國",
	}

	for i, record := range records {
		InsertRecord(db, i+1, record)
	}
	log.Println("Records inserted successfully")

	//LoadFilesystemData(db, "/")
}

func LoadFilesystemData(db *sql.DB, path string) {
	err := ReadDirRecursive(db, path)
	if err != nil {
		log.Fatalf("Error walking the directory: %v", err)
	}
}

func ReadDirRecursive(db *sql.DB, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil && !os.IsPermission(err) {
		return err
	}

	for _, entry := range entries {
		path := dir + "/" + entry.Name()
		if entry.IsDir() {
			// Recursively read subdirectory
			err := ReadDirRecursive(db, path)
			if err != nil {
				return err
			}
		} else {
			//log.Println(entry.Name())
			InsertRecord(db, 0, path)
		}
	}
	return nil
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
