package models

type DroneType string

const (
	DroneMiner     DroneType = "drone_miner"
	DroneTransport DroneType = "drone_transport"
	DroneFighter   DroneType = "drone_fighter"
)

type Drone struct {
	*CoreEntity
	Name string
	Type DroneType
}

func NewDrone(name string, t DroneType) *Drone {
	return &Drone{
		CoreEntity: NewCoreEntity(),
		Name:       name,
		Type:       t,
	}
}
