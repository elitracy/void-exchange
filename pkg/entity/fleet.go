package entity

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
	HP   int
}

func NewDrone(name string, dt DroneType, hp int) *Drone {
	return &Drone{
		CoreEntity: NewCoreEntity(),
		Name:       name,
		Type:       dt,
		HP:         hp,
	}
}

func (d *Drone) Tick() error {
	return nil
}
