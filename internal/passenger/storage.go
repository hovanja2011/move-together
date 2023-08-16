package passenger

import "context"

type Repository interface {
	Create(ctx context.Context, passenger *Passenger) error
	FindOne(ctx context.Context, id string) (Passenger, error)
	FindAll(ctx context.Context) ([]Passenger, error)
	Update(ctx context.Context, passenger Passenger) error
	Delete(ctx context.Context, id string) error
}
