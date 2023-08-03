package ride

import "github.com/hovanja2011/move-together/internal/driver"

type CreateRideDTO struct {
	Driver      driver.Driver `json:"driver"`
	FromAddress string        `json:"from_addr"`
	ToAddress   string        `json:"to_addr"`
}
