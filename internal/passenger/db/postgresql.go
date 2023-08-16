package passenger

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hovanja2011/move-together/internal/passenger"
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

func (r *repository) Create(ctx context.Context, passenger *passenger.Passenger) error {
	q := `
		INSERT 
			INTO passenger (name, score) 
		VALUES 
			($1, $2) 
		RETURNING 
			id`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", formatQuery(q)))
	if err := r.client.QueryRow(ctx, q, passenger.Name, passenger.Score).Scan(&passenger.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr := err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s , Detail: %s , Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return newErr
		}
		return err
	}
	return nil
}

func (r *repository) FindOne(ctx context.Context, id string) (passenger.Passenger, error) {
	q := `
		SELECT 
			id, name, score
		FROM
			public.passenger
		WHERE
			id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", formatQuery(q)))

	var pssg passenger.Passenger
	err := r.client.QueryRow(ctx, q, id).Scan(&pssg.ID, &pssg.Name, &pssg.Score)
	if err != nil {
		return passenger.Passenger{}, err
	}

	return pssg, nil
}

func (r *repository) FindAll(ctx context.Context) ([]passenger.Passenger, error) {
	q := `
		SELECT 
			id, name, score
		FROM
			public.passenger
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	passengers := make([]passenger.Passenger, 0)
	for rows.Next() {
		var pssg passenger.Passenger

		err = rows.Scan(&pssg.ID, &pssg.Name, &pssg.Score)
		if err != nil {
			return nil, err
		}
		passengers = append(passengers, pssg)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return passengers, nil
}

func (r *repository) Update(ctx context.Context, passenger passenger.Passenger) error {
	panic("implement me")
}

func (r *repository) Delete(ctx context.Context, id string) error {
	panic("implement me")
}

func NewRepository(client postgresql.Client, logger *logging.Logger) passenger.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
