package models

import "time"

// StartFragment Fragment interface for "START" fragment.
type StartFragment interface {
	At() time.Time
	Save([]byte) error
	Load() ([]byte, error)
}

// FullFragment Fragment interface for "FULL" fragment.
type FullFragment interface {
	At() time.Time
	Tick() uint64
	Save([]byte) error
	Load() ([]byte, error)
}

// DeltaFragment Fragment interface for "DELTA" fragment.
type DeltaFragment interface {
	At() time.Time
	EndTick() uint64
	Save([]byte) error
	Load() ([]byte, error)
}
