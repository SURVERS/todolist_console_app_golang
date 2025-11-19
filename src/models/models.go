package models

type Task struct {
	ID          int
	Description string
	Completed   bool
}

var Tasks []Task
var NextID int = 1
