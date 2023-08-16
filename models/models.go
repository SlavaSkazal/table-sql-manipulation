package models

import "time"

// User is a struct of human data.
type User struct {
	Name     string
	Birthday time.Time
	Sex      Sex
}

// Sex is a string type, it is a sex of the human, Sex can only contain predefined values.
type Sex string

const (
	Male   Sex = "Male"
	Female Sex = "Female"
)
