package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gitlab.com/slavaskazal1/ptmk/models"
	"gitlab.com/slavaskazal1/ptmk/storage/sqlite"
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
	case 2:
		if len(os.Args) < 7 {
			fmt.Println(makeErr("not enough args"))
			return
		}

		name := os.Args[2] + " " + os.Args[3] + " " + os.Args[4]

		birthday, err := time.Parse("2006-01-02", os.Args[5])
		if err != nil {
			log.Fatal(err)
		}

		var sex models.Sex
		if os.Args[6] == "male" {
			sex = models.Male
		} else if os.Args[6] == "female" {
			sex = models.Female
		} else {
			fmt.Println(makeErr("not valid sex"))
			return
		}

		db, err := createDB()
		if err != nil {
			log.Fatal(err)
		}

		user := models.User{Name: name, Birthday: birthday, Sex: sex}
		err = db.CreateRecord(user)
		if err != nil {
			log.Fatal(err)
		}
	case 1, 3, 4, 5, 6:
		db, err := createDB()
		if err != nil {
			log.Fatal(err)
		}

		switch num {
		case 1:
			err = db.CreateTable()
		case 3:
			err = db.PrintUniqueRecords()
		case 4:
			//err = db.CreateAutoRecords("\"male\"", 100)
			err = db.CreateAutoRecords(models.Male, 100)
		case 5:
			err = db.PrintRecordsByArguments()
		case 6:
			err = db.PrintRecordsByArgumentsIndexed()
		}
		if err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println(makeErr("wrong 1 arg"))
		return
	}
}

func makeErr(strErr string) error {
	return fmt.Errorf("error: %s", strErr)
}

func createDB() (*sqlite.Database, error) {
	return sqlite.New("./storage/ptmk.db")
}
