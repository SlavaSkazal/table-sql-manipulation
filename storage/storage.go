package storage

import (
	"time"
)

type User struct {
	Name     string
	Birthday time.Time
	Sex      string
}

type Storage interface {
	CreateTable() error
	CreateRecord(user User) error
	PrintUniqueRecords() error
	CreateAutoRecords(sex string, count int) error
	PrintRecordsByArguments() error
	PrintRecordsByArgumentsIndexed() error
}
