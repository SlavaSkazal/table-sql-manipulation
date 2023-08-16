package storage

import (
	"gitlab.com/slavaskazal1/ptmk/models"
)

// Storage is a interface which describes data store methods.
type Storage interface {
	// CreateTable creates a table "users" if it hasn't already been created.
	CreateTable() error
	// CreateRecord creates an entry in table "users" with the given data models.User.
	CreateRecord(models.User) error
	// CreateAutoRecords creates automatic records in table "users", of which "n" (int - 2 parameter)
	// with the specified sex (models.Sex - 1 parameter).
	CreateAutoRecords(models.Sex, int) error
	// PrintUniqueRecords prints all lines and the number of full years with given unique parameters in table "users".
	PrintUniqueRecords() error
	// PrintRecordsByArguments prints all lines with given parameters.
	PrintRecordsByArguments() error
	// PrintRecordsByArgumentsIndexed adds indexes for lines with given unique parameters in table "users".
	PrintRecordsByArgumentsIndexed() error
}
