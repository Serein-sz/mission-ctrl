package repository

import (
	"fmt"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	ID          string
	Commiter    string
	Date        string
	Description string
	Repository  string
}

func (t *Task) String() string {
	return fmt.Sprintf("{ id: %s, commiter: %s, date: %s, description: %s, repository: %s }",
		t.ID,
		t.Commiter,
		t.Date,
		t.Description,
		t.Repository,
	)
}

func NewTask(commiter, date, description, repository string) Task {
	id, _ := gonanoid.New()
	return Task{
		ID:          id,
		Commiter:    commiter,
		Date:        date,
		Description: description,
		Repository:  repository,
	}
}
