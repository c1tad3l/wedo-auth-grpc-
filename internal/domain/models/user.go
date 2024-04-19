package models

type User struct {
	Uuid        string
	Email       string
	PassHash    string
	Phone       string
	DateOfBirth string
	Username    string
}
