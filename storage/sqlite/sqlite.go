package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"gitlab.com/slavaskazal1/ptmk/models"
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

// Database is a struct of database sqlite.
type Database struct {
	dbSql *sql.DB
}

// New create new database in directory "path".
func New(path string) (*Database, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("error: %w %w", errors.New("failed open db"), err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error: %w %w", errors.New("failed ping db"), err)
	}
	return &Database{dbSql: db}, nil
}

// CreateTable creates a table "users" in Database if it hasn't already been created.
func (db *Database) CreateTable() error {
	stat, err := db.dbSql.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name VARCHAR(100), birthday DATE, sex VARCHAR(50) CHECK( sex IN ('Male','Female') ))")
	if err != nil {
		return err
	}
	_, err = stat.Exec()
	if err != nil {
		return err
	}
	return nil
}

// CreateRecord creates an entry in table "users" in Database with the given data models.User.
func (db *Database) CreateRecord(user models.User) error {
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
	_, err = stat.Exec(user.Name, user.Birthday, user.Sex)
	if err != nil {
		return err
	}
	return nil
}

// CreateAutoRecords creates 1 million automatic records in table "users" in Database, of which "n" (int - 2 parameter)
// with the specified sex (models.Sex - 1 parametr) and full name starting with "F".
func (db *Database) CreateAutoRecords(s models.Sex, n int) error {
	existTable, err := db.checkTableExist()
	if err != nil {
		return err
	}
	if !existTable {
		return db.makeErrTableNotExist()
	}

	query, err := makeQuery(s, n, 1000000)
	if err != nil {
		return err
	}
	stat, err := db.dbSql.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stat.Exec()
	if err != nil {
		return err
	}
	return nil
}

// PrintUniqueRecords prints all lines and number of full years with a unique name + date + sex from table "users".
func (db *Database) PrintUniqueRecords() error {
	existTable, err := db.checkTableExist()
	if err != nil {
		return err
	}
	if !existTable {
		return db.makeErrTableNotExist()
	}
	rows, err := db.dbSql.Query("SELECT DISTINCT name, birthday, sex FROM users ORDER BY name")
	if err != nil {
		return err
	}
	printRows(rows)
	return nil
}

// PrintRecordsByArguments prints all lines where sex is male, the name starts with "F".
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

	log.Printf("Time has passed: %v\n", time.Now().Sub(startTime))
	return nil
}

// PrintRecordsByArgumentsIndexed adds indexes for lines where sex is male, the name starts with "F" and prints this lines.
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
	_, err = stat.Exec()
	if err != nil {
		return err
	}

	startTime := time.Now()

	rows, err := db.dbSql.Query("SELECT name, birthday, sex FROM users WHERE sex = 'male' AND name like 'F%'")
	if err != nil {
		return err
	}
	printRows(rows)

	log.Printf("Time has passed with indexed: %v\n", time.Now().Sub(startTime))
	return nil
}

// makeQuery creates a row with automatically generated records to insert into table "users".
func makeQuery(sex models.Sex, countFirstLitName int, countAll int) (string, error) {
	if countFirstLitName > countAll {
		return "", errors.New("very match count")
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
		sb.WriteString("\"" + string(sex) + "\"")
		sb.WriteString(")")
		if i < countAll {
			sb.WriteString(", (")
		}
	}
	return sb.String(), nil
}

// randName creates automatically name.
func randName() string {
	return "\"" + lastNames[rand.Intn(len(lastNames))] + " " +
		firstNames[rand.Intn(len(firstNames))] + " " +
		middleNames[rand.Intn(len(middleNames))] + "\""
}

// randName creates automatically name that starts with "F".
func randNameF() string {
	return "\"" + lastNamesF[rand.Intn(20)] + " " +
		firstNames[rand.Intn(20)] + " " +
		middleNames[rand.Intn(20)] + "\""
}

// randSex creates a random value of models.Sex.
func randSex() models.Sex {
	if rand.Intn(2) == 0 {
		return models.Male
	} else {
		return models.Female
	}
}

// randSex creates a random birthday.
func randBirthday() time.Time {
	year := 1900 + rand.Intn(123)
	month := months[1+rand.Intn(11)]
	day := 1 + rand.Intn(28)
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

// printRows prints rows from the query.
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

// checkTableExist checks if a table exists.
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

// makeErrTableNotExist create an error that the table does not exist.
func (db *Database) makeErrTableNotExist() error {
	return errors.New("table not exists")
}
