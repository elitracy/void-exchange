package entity

type DroneType string

const (
	DroneMiner     DroneType = "drone_miner"
	DroneTransport DroneType = "drone_transport"
	DroneFighter   DroneType = "drone_fighter"
)

// DroneActivity tracks what a drone is currently committed to doing.
// Dispatch/recall are instant in Phase 1 (no travel time), so there is no
// separate "en route" or "recalled" transitional state: recall just resets
// a drone straight back to Idle.
type DroneActivity string

const (
	DroneIdle     DroneActivity = "drone_idle"
	DroneFighting DroneActivity = "drone_fighting"
	DroneMining   DroneActivity = "drone_mining"
)

// Level-derived combat stats. Variation by drone type/tier comes later; for
// now level is the single knob for both HP and Attack.
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

// NewDrone creates a drone whose HP and Attack are derived from level (level
// 1 is the baseline). Level below 1 is clamped to 1.
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
