package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresClient struct {
	Client *pgxpool.Pool
}

type Options struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
}

type Person struct {
	PersonID int32
	Name     string
	Age      int32
	Address  string
	Work     string
}

var (
	postgresClient *PostgresClient
	AllFields      = []string{
		"id",
		"name",
		"age",
		"address",
		"work",
	}
)

func NewPostgresClient(ctx context.Context, options *Options) (*PostgresClient, error) {
	if postgresClient == nil {
		client, err := getPostgresClient(ctx, options)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize postgres client, error is: %s", err)
		}

		postgresClient = &PostgresClient{
			Client: client,
		}

		return postgresClient, nil
	}

	return postgresClient, nil
}

func (pc *PostgresClient) CreateTable(ctx context.Context) error {
	_, err := pc.Client.Query(
		ctx,
		"CREATE TABLE persons (id serial PRIMARY KEY,  name VARCHAR (64), age INT, address VARCHAR(128), work VARCHAR(128));",
	)

	if err != nil {
		return fmt.Errorf("failed to run postgres query, error is: %s", err)
	}

	return nil
}

func (pc *PostgresClient) Create(ctx context.Context, table string, fields []string, person Person) (int32, error) {
	rows, err := pc.Client.Query(
		ctx,
		fmt.Sprintf(
			"INSERT INTO %s(%s, %s, %s, %s) VALUES ('%s', %d, '%s', '%s') RETURNING id",
			table,
			fields[1],
			fields[2],
			fields[3],
			fields[4],
			person.Name,
			person.Age,
			person.Address,
			person.Work,
		),
	)

	if err != nil {
		return 0, fmt.Errorf("failed to run postgres query, error is: %s", err)
	}

	var personId int32
	for rows.Next() {
		if err := rows.Scan(&personId); err != nil {
			return personId, fmt.Errorf("failed to map row to data, error is: %s", err)
		}
	}

	return personId, nil
}

func (pc *PostgresClient) Delete(ctx context.Context, table string, personID int) error {
	_, err := pc.Client.Query(
		ctx,
		fmt.Sprintf(
			"DELETE FROM %s WHERE id = '%d'",
			table,
			personID,
		),
	)

	if err != nil {
		return fmt.Errorf("failed to run postgres query, error is: %s", err)
	}

	return nil
}

func (pc *PostgresClient) Get(ctx context.Context, table string, personID int) ([]Person, error) {
	rows, err := pc.Client.Query(
		ctx,
		fmt.Sprintf(
			"SELECT * FROM %s WHERE id = %d",
			table,
			personID,
		),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to run postgres query, error is: %s", err)
	}

	data := make([]Person, 0)

	for rows.Next() {
		person := Person{}
		err := rows.Scan(
			&person.PersonID,
			&person.Name,
			&person.Age,
			&person.Address,
			&person.Work,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to map row to data, error is: %s", err)
		}

		data = append(data, person)
	}

	return data, nil
}

func (pc *PostgresClient) GetAll(ctx context.Context, table string) ([]Person, error) {
	rows, err := pc.Client.Query(
		ctx,
		fmt.Sprintf(
			"SELECT * FROM %s",
			table,
		),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to run postgres query, error is: %s", err)
	}

	data := make([]Person, 0)

	for rows.Next() {
		person := Person{}
		err := rows.Scan(
			&person.PersonID,
			&person.Name,
			&person.Age,
			&person.Address,
			&person.Work,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to map row to data, error is: %s", err)
		}

		data = append(data, person)
	}

	return data, nil
}

func (pc *PostgresClient) Update(ctx context.Context, table string, fields []string, person Person) error {
	if person.Name != "" {
		_, err := pc.Client.Query(
			ctx,
			fmt.Sprintf(
				"UPDATE %s SET %s = '%s' WHERE id = '%d'",
				table,
				fields[1],
				person.Name,
				person.PersonID,
			),
		)

		if err != nil {
			return fmt.Errorf("failed to run postgres query, error is: %s", err)
		}
	}

	if person.Age != 0 {
		_, err := pc.Client.Query(
			ctx,
			fmt.Sprintf(
				"UPDATE %s SET %s = '%d' WHERE id = '%d'",
				table,
				fields[2],
				person.Age,
				person.PersonID,
			),
		)

		if err != nil {
			return fmt.Errorf("failed to run postgres query, error is: %s", err)
		}
	}

	if person.Address != "" {
		_, err := pc.Client.Query(
			ctx,
			fmt.Sprintf(
				"UPDATE %s SET %s = '%s' WHERE id = '%d'",
				table,
				fields[3],
				person.Address,
				person.PersonID,
			),
		)

		if err != nil {
			return fmt.Errorf("failed to run postgres query, error is: %s", err)
		}
	}

	if person.Work != "" {
		_, err := pc.Client.Query(
			ctx,
			fmt.Sprintf(
				"UPDATE %s SET %s = '%s' WHERE id = '%d'",
				table,
				fields[4],
				person.Work,
				person.PersonID,
			),
		)

		if err != nil {
			return fmt.Errorf("failed to run postgres query, error is: %s", err)
		}
	}

	return nil
}

func getPostgresClient(ctx context.Context, options *Options) (*pgxpool.Pool, error) {
	dbURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		options.Host, options.Port, options.User, options.DB, options.Password)

	client, err := pgxpool.Connect(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect postgres client, error is: %s", err)
	}

	return client, nil
}
