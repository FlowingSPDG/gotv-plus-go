package gotv

import (
	"time"
)

// Store Store interface is call-back interface for Storing GOTV+ fragments. WRITE ONLY OPERATION.
type Store interface {
	Auth(token string, auth string) error
	OnStart(token string, fragment int, f StartFrame) error
	OnFull(token string, fragment int, tick int, at time.Time, b []byte) error
	OnDelta(token string, fragment int, endtick int, at time.Time, final bool, b []byte) error
}

// Broadcaster GOTV+ broadcasts GOTV+ demo fragments to CS:GO playcast clients. READ ONLY OPERATION.
type Broadcaster interface {
	GetSync(token string, fragment int) (Sync, error)
	GetSyncLatest(token string) (Sync, error)
	GetStart(token string, fragment int) ([]byte, error) // could be io.Reader??
	GetFull(token string, fragment int) ([]byte, error)
	GetDelta(token string, fragment int) ([]byte, error)
}

// Fragment has both of Full/Delta fragment data
type Fragment struct {
	At      time.Time
	Tick    int
	Final   bool
	EndTick int
	Full    []byte
	Delta   []byte
}

// StartFrame Start fragment
type StartFrame struct {
	At       time.Time
	Tps      float64 // Even though it is int, we should use float64 because server sends its value as "128.0"
	Protocol int
	Map      string
	Body     []byte
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

// StartQuery Query for START request
type StartQuery struct {
	Tick             int     `query:"tick" form:"tick"` // the starting tick of the broadcast
	TPS              float64 `query:"tps" form:"tps"`   // the tickrate of the GOTV broadcast. // 実際はintだが128.0 という小数点付きで送られてくるのでfloatに設定する
	Map              string  `query:"map" form:"map"`   // the name of the map
	KeyframeInterval float64 `json:"keyframe_interval,omitempty"`
	Protocol         int     `query:"protocol" form:"protocol"` // Currently 4
}

// FullQuery Query for FULL request
type FullQuery struct {
	Tick int `query:"tick" form:"tick"` // the starting tick of the broadcast
}

// DeltaQuery Query for DELTA request
type DeltaQuery struct {
	EndTick int  `query:"endtick" form:"endtick"` // endtick of delta frame
	Final   bool `query:"final" form:"final"`     // is final fragment
}

// SyncQuery Query for SYNC request
type SyncQuery struct {
	Fragment int `query:"fragment" form:"fragment"` // endtick of delta frame
}
