package domain

import "time"

type ExampleItem struct {
	ID          string
	Name        string
	Description string
	CreateAt    time.Time
	UpdateAt    time.Time
	DeletedAt   time.Time
}
