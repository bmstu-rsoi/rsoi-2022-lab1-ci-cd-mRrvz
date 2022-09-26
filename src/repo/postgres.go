package repo

import (
	"context"
	"fmt"

	"github.com/bmstu-rsoi/rsoi-2022-lab1-ci-cd-mRrvz/src/db"
)

type postgresRepo struct {
	Client db.PostgresClient
}

func NewPostgresRepo(client db.PostgresClient) PostgresRepository {
	return &postgresRepo{
		Client: client,
	}
}

func (pr *postgresRepo) CreateTable() error {
	if err := pr.Client.CreateTable(context.Background()); err != nil {
		return err
	}

	return nil
}

func (pr *postgresRepo) Create(person db.Person) (int32, error) {
	personId, err := pr.Client.Create(context.Background(), "persons", db.AllFields, person)
	if err != nil {
		return personId, fmt.Errorf("failed to insert person, error is: %s", err)
	}

	return personId, nil
}

func (pr *postgresRepo) Get(personID int) (*db.Person, error) {
	persons, err := pr.Client.Get(context.Background(), "persons", personID)
	if err != nil {
		return nil, fmt.Errorf("failed to get person with id = %d, error is: %s", personID, err)
	}

	if len(persons) == 0 {
		return nil, nil
	}

	return &persons[0], nil
}

func (pr *postgresRepo) GetAll() ([]db.Person, error) {
	persons, err := pr.Client.GetAll(context.Background(), "persons")
	if err != nil {
		return nil, fmt.Errorf("failed to get all persons, error is: %s", err)
	}

	return persons, nil
}

func (pr *postgresRepo) Update(person db.Person, personID int) error {
	err := pr.Client.Update(context.Background(), "persons", db.AllFields, person)
	if err != nil {
		return fmt.Errorf("failed to insert person, error is: %s", err)
	}

	return nil
}

func (pr *postgresRepo) Delete(personID int) error {
	err := pr.Client.Delete(context.Background(), "persons", personID)
	if err != nil {
		return fmt.Errorf("failed to delete person with id = %d, error is: %s", personID, err)
	}

	return nil
}
