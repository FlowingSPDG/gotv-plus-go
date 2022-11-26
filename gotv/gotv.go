package gotv

import (
	"time"
)

const (
	authHeader = "X-Origin-Auth"
)

// Store Store interface is call-back interface for Storing GOTV+ fragments. WRITE ONLY OPERATION.
type Store interface {
	Auth(token string, auth string) error
	OnStart(token string, f StartFrame) error
	OnFull(token string, f FullFrame) error
	OnDelta(token string, f DeltaFrame) error
}

// Broadcaster GOTV+ broadcasts GOTV+ demo fragments to CS:GO playcast clients. READ ONLY OPERATION.
type Broadcaster interface {
	GetSync(token string) error
	GetStart(token string, fragment int) ([]byte, error) // could be io.Reader??
	GetFull(token string, fragment int, tick int) ([]byte, error)
	GetDelta(token string, fragment int, endtick int) ([]byte, error)
}

// StartFrame Start fragment
type StartFrame struct {
	At       time.Time
	Tps      int
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
