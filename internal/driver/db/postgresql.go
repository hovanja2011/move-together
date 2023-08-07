package driver

import (
	"context"
	"fmt"
	"strings"

	"github.com/hovanja2011/move-together/internal/driver"
	"github.com/hovanja2011/move-together/pkg/client/postgresql"
	"github.com/hovanja2011/move-together/pkg/logging"
	"github.com/jackc/pgconn"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", " "), "\n", " ")
}

func (r *repository) Create(ctx context.Context, driver *driver.Driver) error {
	q := `
		INSERT 
			INTO driver (name, score) 
		VALUES 
			($1, $2) 
		RETURNING 
			id`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", formatQuery(q)))
	if err := r.client.QueryRow(ctx, q, driver.Name, driver.Score).Scan(&driver.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s , Detail: %s , Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return newErr
		}
		return err
	}
	return nil
}

func (r *repository) FindOne(ctx context.Context, id string) (driver.Driver, error) {
	q := `
		SELECT 
			id, name, score
		FROM
			public.driver
		WHERE
			id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", formatQuery(q)))

	var drv driver.Driver
	err := r.client.QueryRow(ctx, q, id).Scan(&drv.ID, &drv.Name, &drv.Score)
	if err != nil {
		return driver.Driver{}, err
	}

	return drv, nil
}

func (r *repository) FindAll(ctx context.Context) ([]driver.Driver, error) {
	q := `
		SELECT 
			id, name, score
		FROM
			public.driver
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	drivers := make([]driver.Driver, 0)
	for rows.Next() {
		var drv driver.Driver

		err = rows.Scan(&drv.ID, &drv.Name, &drv.Score)
		if err != nil {
			return nil, err
		}
		drivers = append(drivers, drv)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return drivers, nil
}

func (r *repository) Update(ctx context.Context, driver driver.Driver) error {
	panic("implement me")
}

func (r *repository) Delete(ctx context.Context, id string) error {
	panic("implement me")
}

func NewRepository(client postgresql.Client, logger *logging.Logger) driver.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
