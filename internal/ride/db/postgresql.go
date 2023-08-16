package ride

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hovanja2011/move-together/internal/ride"
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

func (r *repository) Create(ctx context.Context, ride *ride.Ride) error {
	q := `
		INSERT 
			INTO ride (driver_id, from_addr, to_addr) 
		VALUES 
			($1, $2, $3) 
		RETURNING 
			id`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", formatQuery(q)))
	if err := r.client.QueryRow(ctx, q, ride.DriverId, ride.FromAddress, ride.ToAddress).Scan(&ride.ID); err != nil {
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

func (r *repository) FindOne(ctx context.Context, id string) (ride.Ride, error) {
	q := `
		SELECT 
			id, driver_id, from_addr, to_addr
		FROM
			public.ride
		WHERE
			id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", formatQuery(q)))

	var rd ride.Ride
	err := r.client.QueryRow(ctx, q, id).Scan(&rd.ID, &rd.DriverId, &rd.FromAddress, &rd.ToAddress)
	if err != nil {
		return ride.Ride{}, err
	}

	return rd, nil
}

func (r *repository) FindAll(ctx context.Context) ([]ride.Ride, error) {
	q := `
		SELECT 
			id, driver_id, from_addr, to_addr
		FROM
			public.ride
	`
	r.logger.Trace(fmt.Sprintf("SQL query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	rides := make([]ride.Ride, 0)
	for rows.Next() {
		var rd ride.Ride

		err = rows.Scan(&rd.ID, &rd.DriverId, &rd.FromAddress, &rd.ToAddress)
		if err != nil {
			return nil, err
		}
		rides = append(rides, rd)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return rides, nil
}

func (r *repository) Update(ctx context.Context, ride ride.Ride) error {
	panic("implement me")
}

func (r *repository) Delete(ctx context.Context, id string) error {
	panic("implement me")
}

func NewRepository(client postgresql.Client, logger *logging.Logger) ride.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
