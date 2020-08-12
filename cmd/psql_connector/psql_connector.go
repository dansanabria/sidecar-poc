package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
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

func fileExistCheck(filepath string) bool {
	fileinfo, err := os.Stat(filepath)
	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Printf("File size is %s\n", strconv.FormatInt(fileinfo.Size(), 10))
	return true
}

func fileSizeCheck(filepath string) int64 {
	fileinfo, err := os.Stat(filepath)
	if err != nil {
		log.Fatal(err)
	}

	return fileinfo.Size()
}

func main() {
	host := os.Getenv("dbhost")
	port := 5432
	user := os.Getenv("dbuser")
	dbname := os.Getenv("dbname")

	for fileExistCheck("/token/.token") != true {
		fmt.Println("Token does not exist, sleeping for 5 secs...")
		time.Sleep(5 * time.Second)
	}

	for fileSizeCheck("/token/.token") == 0 {
		fmt.Println("Token being retrieved. Sleeping for 5 secs...")
		time.Sleep(5 * time.Second)
	}

	password := readToken("/token/.token")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=require",
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
