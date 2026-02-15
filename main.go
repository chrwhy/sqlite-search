package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/chrwhy/simple/examples/go/im-search"
	load_test "github.com/chrwhy/simple/examples/go/load-test"
	"github.com/chrwhy/simple/examples/go/qparser"
	"github.com/chrwhy/simple/examples/go/util"
)

func ImSearchInit(db *sql.DB) {
	im_search.CreateChatGroupTable(db)
	im_search.CreateGroupMemberTable(db)
	im_search.CreateChatMessageTable(db)
	im_search.CreateContactTable(db)

	// Seed example chat groups (no-op if already seeded).
	im_search.SeedChatGroups(db)
	// Seed example group members (no-op if already seeded).
	im_search.SeedGroupMembers(db)
	// Seed example contacts (no-op if already seeded).
	im_search.SeedContacts(db)
	// Seed example chat messages (no-op if already seeded).
	im_search.SeedChatMessages(db)
}

func ExternalSearchInit(db *sql.DB) {
	load_test.CreateExternalTable(db)
	load_test.MockData(db)
}

func main() {
	db := util.InitDB()
	defer db.Close()

	//ImSearchInit(db)
	//ExternalSearchInit(db)
	//spotlight.InitData(db)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Simple SQL REPL")
	fmt.Println("---------------------")

	for {
		fmt.Println("Choose an option:")
		fmt.Println("1. Insert record")
		fmt.Println("2. Query Mode")
		fmt.Println("3. SQL Mode")
		fmt.Println("4. Exit")
		fmt.Print("Enter choice: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			for {
				fmt.Print("Enter text to insert (type 'exit' to go back): ")
				text, _ := reader.ReadString('\n')
				text = strings.TrimSpace(text)
				if len(text) == 0 {
					continue
				}
				if strings.ToLower(text) == "exit" {
					break
				}
				InsertRecord(db, 0, text)
			}
		case "2":
			for {
				fmt.Print("Enter Query Mode (type 'exit' to go back): ")
				query, _ := reader.ReadString('\n')
				query = strings.TrimSpace(query)
				if len(query) == 0 {
					continue
				}
				if strings.ToLower(query) == "exit" {
					break
				}

				clause := qparser.ParseClause(query)
				sql := "select biz_id, simple_highlight(t1, 1, '[', ']') from t1 where text match ('" + clause + "')"
				log.Println(sql)
				util.Query(db, sql)

				log.Println(strings.Repeat("=", 60))
				log.Println(">>>>>>>>>> Chat Groups")
				a, _ := im_search.SearchChatGroups(db, clause)
				if len(a) != 0 {
					log.Println(a)
				} else {
					log.Println("(no results)")
				}
				log.Println()

				log.Println(strings.Repeat("=", 60))
				log.Println(">>>>>>>>>> Chat Messages")
				b, _ := im_search.SearchChatMessages(db, query)
				if len(b) != 0 {
					log.Println(b)
				} else {
					log.Println("(no results)")
				}
				log.Println()

				log.Println(strings.Repeat("=", 60))
				log.Println(">>>>>>>>>> Contacts")
				c, _ := im_search.SearchContacts(db, clause)
				if len(c) != 0 {
					log.Println(c)
				} else {
					log.Println("(no results)")
				}
				log.Println()

				log.Println(strings.Repeat("=", 60))
				log.Println(">>>>>>>>>> Group Members")
				d, _ := im_search.SearchGroupMembers(db, clause)
				if len(d) != 0 {
					log.Println(d)
				} else {
					log.Println("(no results)")
				}
				log.Println(strings.Repeat("=", 60))
			}
		case "3":
			for {
				fmt.Print("Enter SQL Mode (type 'exit' to go back): ")
				query, _ := reader.ReadString('\n')
				query = strings.TrimSpace(query)
				if len(query) == 0 {
					continue
				}
				if strings.ToLower(query) == "exit" {
					break
				}
				log.Println(query)
				util.Query(db, query)
			}
		case "4":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
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
