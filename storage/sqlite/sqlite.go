package sqlite

import (
	"database/sql"
	"fmt"
	"gitlab.com/slavaskazal1/ptmk/storage"
	"math/rand"
	"strings"
	"time"
)

var firstNames = [20]string{"Aleksandr", "Aleksey", "Andrey", "Anton", "Artem", "Boris", "Vadim", "Vasily", "Vladimir",
	"Georgy", "Denis", "Dmitry", "Egor", "Ivan", "Maksim", "Nikita", "Oleg", "Pavel", "Petr", "Roman"}

var middleNames = [20]string{"Aleksandrovich", "Alexeyevich", "Andreevich", "Antonovich", "Artemovich", "Borisovich", "Vadimovich", "Vasilevich", "Vladimirovich",
	"Georgievich", "Denisovich", "Dmitrievich", "Egorovich", "Ivanovich", "Maksimovich", "Nikitich", "Olegovich", "Pavlovich", "Petrovich", "Romanovich"}

var lastNames = [20]string{"Asin", "Enin", "Svetin", "Antonin", "Artemin", "Borisov", "Vadimov", "Vasin", "Vladin",
	"Georgiev", "Denisov", "Dmitriev", "Egorov", "Ivanov", "Maksimov", "Nikitin", "Olgin", "Pavlov", "Petrov", "Romanin"}

var lastNamesF = [20]string{"Fasin", "Fenin", "Fetlin", "Fantonin", "Femin", "Fisov", "Fadimov", "Flasin", "Filadin",
	"Fergiev", "Fenisov", "Fimitriev", "Forov", "Fanov", "Flagov", "Frakitin", "Folgin", "Fikin", "Fetrov", "Fomanin"}

var months = [12]time.Month{time.January, time.February, time.March, time.April, time.May, time.June,
	time.July, time.August, time.September, time.October, time.November, time.December}

type Database struct {
	dbSql *sql.DB
}

func New(path string) (*Database, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("error open db: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error ping db: %w", err)
	}
	return &Database{dbSql: db}, nil
}

func (db *Database) CreateTable() error {
	stat, err := db.dbSql.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT, birthday DATE, sex TEXT)")
	if err != nil {
		return err
	}
	stat.Exec()
	return nil
}

func (db *Database) CreateRecord(user storage.User) error {
	existTable, err := db.checkTableExist()
	if err != nil {
		return err
	}
	if !existTable {
		return db.makeErrTableNotExist()
	}

	stat, err := db.dbSql.Prepare("INSERT INTO users (name, birthday, sex) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	stat.Exec(user.Name, user.Birthday, user.Sex)
	return nil
}

func (db *Database) PrintUniqueRecords() error {
	existTable, err := db.checkTableExist()
	if err != nil {
		return err
	}
	if !existTable {
		return db.makeErrTableNotExist()
	}
	//надо уникальные только по имя+дата, но не по полу
	rows, err := db.dbSql.Query("SELECT DISTINCT name, birthday, sex FROM users ORDER BY name") //WHERE
	if err != nil {
		return err
	}
	printRows(rows)
	return nil
}

func (db *Database) CreateAutoRecords(sex string, count int) error {
	existTable, err := db.checkTableExist()
	if err != nil {
		return err
	}
	if !existTable {
		return db.makeErrTableNotExist()
	}

	query, err := makeQuery(sex, count, 1000000)
	if err != nil {
		return err
	}
	stat, err := db.dbSql.Prepare(query)
	if err != nil {
		return err
	}
	stat.Exec()
	return nil
}

func (db *Database) PrintRecordsByArguments() error {
	existTable, err := db.checkTableExist()
	if err != nil {
		return err
	}
	if !existTable {
		return db.makeErrTableNotExist()
	}

	startTime := time.Now()

	rows, err := db.dbSql.Query("SELECT name, birthday, sex FROM users WHERE sex = 'male' AND name like 'F%'")
	if err != nil {
		return err
	}
	printRows(rows)

	fmt.Printf("Time has passed: %v\n", time.Now().Sub(startTime))
	return nil
}

func (db *Database) PrintRecordsByArgumentsIndexed() error {
	existTable, err := db.checkTableExist()
	if err != nil {
		return err
	}
	if !existTable {
		return db.makeErrTableNotExist()
	}

	stat, err := db.dbSql.Prepare("CREATE INDEX IF NOT EXISTS idx_male_name_f ON users (sex, name) WHERE sex = 'male'  AND name like 'F%'")
	if err != nil {
		return err
	}
	stat.Exec()

	startTime := time.Now()

	rows, err := db.dbSql.Query("SELECT name, birthday, sex FROM users WHERE sex = 'male' AND name like 'F%'")
	if err != nil {
		return err
	}
	printRows(rows)

	fmt.Printf("Time has passed with indexed: %v\n", time.Now().Sub(startTime))
	return nil
}

func makeQuery(sex string, countFirstLitName int, countAll int) (string, error) {
	if countFirstLitName > countAll {
		return "", fmt.Errorf("error: %v", "very match count")
	}
	sb := strings.Builder{}
	sb.WriteString("INSERT INTO users (name, birthday, sex) VALUES (")
	var name string
	for i := 1; i <= countAll; i++ {
		if countFirstLitName == 0 {
			name = randName()
			sex = randSex()
		} else {
			name = randNameF()
			countFirstLitName--
		}
		birthday := "'" + randBirthday().String() + "'"
		sb.WriteString(name)
		sb.WriteString(", ")
		sb.WriteString(birthday)
		sb.WriteString(", ")
		sb.WriteString(sex)
		sb.WriteString(")")
		if i < countAll {
			sb.WriteString(", (")
		}
	}
	return sb.String(), nil
}

func randName() string {
	return "\"" + lastNames[rand.Intn(20)] + " " + firstNames[rand.Intn(20)] + " " + middleNames[rand.Intn(20)] + "\""
}

func randNameF() string {
	return "\"" + lastNamesF[rand.Intn(20)] + " " + firstNames[rand.Intn(20)] + " " + middleNames[rand.Intn(20)] + "\""
}

func randSex() string {
	if rand.Intn(2) == 0 {
		return "\"male\""
	} else {
		return "\"female\""
	}
}

func randBirthday() time.Time {
	year := 1900 + rand.Intn(123)
	month := months[1+rand.Intn(11)]
	day := 1 + rand.Intn(28)
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func printRows(rows *sql.Rows) {
	var name, sex string
	var birthday time.Time
	now := time.Now()
	for rows.Next() {
		rows.Scan(&name, &birthday, &sex)
		age := int(now.Sub(birthday)/time.Hour) / 8760
		fmt.Println(name, birthday, sex, age)
	}
}

func (db *Database) checkTableExist() (bool, error) {
	rows, err := db.dbSql.Query("SELECT COUNT (*) FROM sqlite_master WHERE type='table' AND name='users'")
	if err != nil {
		return false, err
	}
	var exist int
	for rows.Next() {
		rows.Scan(&exist)
	}
	if exist == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func (db *Database) makeErrTableNotExist() error {
	return fmt.Errorf("error: %v", "table not exist")
}
