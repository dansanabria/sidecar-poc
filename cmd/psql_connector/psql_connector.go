package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func readToken(tokenPath string) string {
	file, err := os.Open(tokenPath)
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	return string(data)
}

func main() {
	host := os.Getenv("dbhost")
	port := 5432
	user := os.Getenv("dbuser")
	dbname := os.Getenv("dbname")
	password := readToken("/token/.token")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	fmt.Println("Sleeping for 5 min ...")
	time.Sleep(300 * time.Second)
}
