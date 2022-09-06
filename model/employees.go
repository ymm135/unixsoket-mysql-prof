package model

import "time"

type Employees struct {
	Id        int
	BirthDate time.Time `gorm:"column:birth_date"`
	FirstName string    `gorm:"column:first_name"`
	LastName  string    `gorm:"column:last_name"`
	Gender    string    `gorm:"column:gender"`
	HireDate  time.Time `gorm:"column:hire_date"`
}
