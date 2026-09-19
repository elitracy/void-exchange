package gamestate

import "errors"

var (
	ErrDroneNotOwned     = errors.New("drone does not belong to faction")
	ErrDroneNotIdle      = errors.New("drone is not idle")
	ErrDroneNotCommitted = errors.New("drone is not currently dispatched")
	ErrWrongDroneType    = errors.New("drone type cannot perform this activity")
	ErrTerritoryNotOwned = errors.New("faction does not own this territory")
	ErrInvalidActivity   = errors.New("invalid dispatch activity")
	ErrNoDronesSpecified = errors.New("no drones specified")
)
