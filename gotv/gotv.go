package gotv

import (
	"time"
)

// Store Store interface is call-back interface for Storing GOTV+ fragments. WRITE ONLY OPERATION.
type Store interface {
	Auth(token string, auth string) error
	OnStart(token string, fragment int, f StartFrame) error
	OnFull(token string, fragment int, f FullFrame) error
	OnDelta(token string, fragment int, f DeltaFrame) error
}

// Broadcaster GOTV+ broadcasts GOTV+ demo fragments to CS:GO playcast clients. READ ONLY OPERATION.
type Broadcaster interface {
	GetSync(token string) (Sync, error)
	GetStart(token string, fragment int) ([]byte, error) // could be io.Reader??
	GetFull(token string, fragment int) ([]byte, error)
	GetDelta(token string, fragment int) ([]byte, error)
}

// StartFrame Start fragment
type StartFrame struct {
	At       time.Time
	Tps      float64 // Even though it is int, we should use float64 because server sends its value as "128.0"
	Protocol int
	Map      string
	Body     []byte
}

// FullFrame Full Fragment
type FullFrame struct {
	At   time.Time
	Tick int
	Body []byte
}

// DeltaFrame Delta fragment
type DeltaFrame struct {
	At      time.Time
	Final   bool
	EndTick int
	Body    []byte
}

// Sync /sync JSON
type Sync struct {
	Tick             int     `json:"tick"`
	Endtick          int     `json:"endtick,omitempty"`
	RealTimeDelay    float64 `json:"rtdelay,omitempty"`
	ReceiveAge       float64 `json:"rcvage,omitempty"`
	Fragment         int     `json:"fragment"`
	SignupFragment   int     `json:"signup_fragment"`
	TickPerSecond    int     `json:"tps"`
	KeyframeInterval float64 `json:"keyframe_interval,omitempty"`
	TokenRedirect    string  `json:"token_redirect,omitempty"`
	Map              string  `json:"map"`
	Protocol         int     `json:"protocol"`
}
