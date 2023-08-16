package ride

type Ride struct {
	ID          string `json:"id"`
	DriverId    string `json:"driver_id"`
	FromAddress string `json:"from_addr"`
	ToAddress   string `json:"to_addr"`
}
