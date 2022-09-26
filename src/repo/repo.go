package repo

import "github.com/bmstu-rsoi/rsoi-2022-lab1-ci-cd-mRrvz/src/db"

type PostgresRepository interface {
	CreateTable() error
	Create(person db.Person) (int32, error)
	Get(userID int) (*db.Person, error)
	GetAll() ([]db.Person, error)
	Update(preson db.Person, userID int) error
	Delete(userID int) error
}
