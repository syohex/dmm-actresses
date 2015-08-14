package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/syohex/dmm/actress"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: actresses.go DB_file")
		os.Exit(1)
	}

	db, err := sql.Open("sqlite3", os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`DELETE FROM actresses`)
	if err != nil {
		log.Println(err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("INSERT INTO actresses(id, name, image) values(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	keywords := actress.Keywords()
	for _, keyword := range keywords {
		log.Printf("Get '%s' actress\n", keyword)

		actresses, err := actress.CollectFromKey(keyword)
		if err != nil {
			log.Fatalln(err)
		}

		for _, actress := range actresses {
			_, err := stmt.Exec(actress.ID, actress.Name, actress.Image)
			if err != nil {
				log.Fatalln(err)
			}
			log.Printf("Insert %v\n", actress)
		}
	}

	tx.Commit()
}
