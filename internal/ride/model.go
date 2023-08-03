package ride

import "github.com/hovanja2011/move-together/internal/driver"

type Ride struct {
	ID          string        `json:"id"`
	Driver      driver.Driver `json:"driver"`
	FromAddress string        `json:"from_addr"`
	ToAddress   string        `json:"to_addr"`
}
