package model

import "time"

type Task struct {
	ID        int
	Title     string
	CreatedAt time.Time
}
