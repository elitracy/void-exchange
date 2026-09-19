package entity

type DroneType string

const (
	DroneMiner     DroneType = "drone_miner"
	DroneTransport DroneType = "drone_transport"
	DroneFighter   DroneType = "drone_fighter"
)

// DroneActivity has no "en route" state: dispatch/recall are instant.
type DroneActivity string

const (
	DroneIdle     DroneActivity = "drone_idle"
	DroneFighting DroneActivity = "drone_fighting"
	DroneMining   DroneActivity = "drone_mining"
)

const (
	DroneBaseHP         = 50
	DroneHPPerLevel     = 25
	DroneBaseAttack     = 5
	DroneAttackPerLevel = 5
)

type Drone struct {
	*CoreEntity
	Name     string
	Type     DroneType
	Level    int
	HP       int
	Attack   int
	Activity DroneActivity
	Target   EntityId // territory currently assigned to; -1 when idle
}

// NewDrone clamps level below 1 up to 1.
func NewDrone(name string, dt DroneType, level int) *Drone {
	if level < 1 {
		level = 1
	}

	return &Drone{
		CoreEntity: NewCoreEntity(),
		Name:       name,
		Type:       dt,
		Level:      level,
		HP:         DroneBaseHP + (level-1)*DroneHPPerLevel,
		Attack:     DroneBaseAttack + (level-1)*DroneAttackPerLevel,
		Activity:   DroneIdle,
		Target:     -1,
	}
}

func (d *Drone) Tick() error {
	return nil
}
