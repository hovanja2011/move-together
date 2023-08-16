package passenger

import (
	"context"

	"github.com/hovanja2011/move-together/pkg/logging"
)

type Service struct {
	repository Repository
	logger     *logging.Logger
}

func (s *Service) Create(ctx context.Context, dto CreatePassengerDTO) (Passenger, error) {
	return Passenger{}, nil
}
