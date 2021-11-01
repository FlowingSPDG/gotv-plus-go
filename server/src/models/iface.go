package models

import "time"

// FragType Fragment type. Start(1), Full(2), or Delta(3).
type FragType int

const (
	// FragTypeStart Fragment for "START" request.
	FragTypeStart = iota + 1
	// FragTypeFull Fragment for "FULL" request.
	FragTypeFull
	// FragTypeDelta Delta Fragment for "DELTA" request.
	FragTypeDelta
)

// Fragment Fragment interface.
type Fragment interface {
	Type() FragType
	At() time.Time
	Tick() uint64    // For full fragment
	EndTick() uint64 // For delta fragment

	Save([]byte) error
	Load() ([]byte, error)
}
