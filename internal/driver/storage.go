package driver

import "context"

type Repository interface {
	Create(ctx context.Context, driver *Driver) error
	FindOne(ctx context.Context, id string) (Driver, error)
	FindAll(ctx context.Context) ([]Driver, error)
	Update(ctx context.Context, driver Driver) error
	Delete(ctx context.Context, id string) error
}
