package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gitlab.com/slavaskazal1/ptmk/storage"
	"gitlab.com/slavaskazal1/ptmk/storage/sqlite"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(makeErr("missing args"))
		return
	}

	num, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(makeErr("not valid 1 arg"))
		return
	}

	switch num {
	case 1:
		db, err := createDB()
		if err != nil {
			fmt.Println(err)
			return
		}
		err = db.CreateTable()
		if err != nil {
			fmt.Println(err)
			return
		}
	case 2:
		//go run main.go 2 Kirin Alexandr Fedorovich 1990-11-10 male
		//go run main.go 2 Olgina Alexandra Igorevna 1995-06-10 female
		if len(os.Args) < 7 {
			fmt.Println(makeErr("not enough args"))
			return
		}
		name := os.Args[2] + " " + os.Args[3] + " " + os.Args[4]
		birthday, err := time.Parse("2006-01-02", os.Args[5])
		if err != nil {
			fmt.Println(err)
			return
		}
		var sex string
		if os.Args[6] == "male" {
			sex = "male"
		} else if os.Args[6] == "female" {
			sex = "female"
		} else {
			fmt.Println(makeErr("not valid sex"))
			return
		}

		db, err := createDB()
		if err != nil {
			fmt.Println(err)
			return
		}
		user := storage.User{Name: name, Birthday: birthday, Sex: sex}
		err = db.CreateRecord(user)
		if err != nil {
			fmt.Println(err)
			return
		}
	case 3:
		db, err := createDB()
		if err != nil {
			fmt.Println(err)
			return
		}
		err = db.PrintUniqueRecords()
		if err != nil {
			fmt.Println(err)
			return
		}
	case 4:
		db, err := createDB()
		if err != nil {
			fmt.Println(err)
			return
		}
		err = db.CreateAutoRecords("\"male\"", 100)
		if err != nil {
			fmt.Println(err)
			return
		}
	case 5:
		db, err := createDB()
		if err != nil {
			fmt.Println(err)
			return
		}
		err = db.PrintRecordsByArguments()
		if err != nil {
			fmt.Println(err)
			return
		}
	case 6:
		db, err := createDB()
		if err != nil {
			fmt.Println(err)
			return
		}
		err = db.PrintRecordsByArgumentsIndexed()
		if err != nil {
			fmt.Println(err)
			return
		}
	default:
		fmt.Println(makeErr("wrong 1 arg"))
		return
	}
}

func makeErr(strErr string) error {
	return fmt.Errorf("error: %v", strErr)
}

func createDB() (*sqlite.Database, error) {
	return sqlite.New("./storage/ptmk.db")
}
