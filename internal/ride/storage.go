package ride

import "context"

type Repository interface {
	Create(ctx context.Context, ride *Ride) error
	FindOne(ctx context.Context, id string) (Ride, error)
	FindAll(ctx context.Context) ([]Ride, error)
	Update(ctx context.Context, ride Ride) error
	Delete(ctx context.Context, id string) error
}
