package models

import "time"

type Startframe struct {
	At   time.Time
	Body []byte
}

type Fullframe struct {
	At   time.Time
	Tick uint64
	Body []byte
}

type Deltaframes struct {
	At      time.Time
	Body    []byte
	EndTick uint64
}

type SyncJSON struct {
	Tick             uint64  `json:"tick"`
	Endtick          uint64  `json:"endtick,omitempty"`
	RealTimeDelay    float64 `json:"rtdelay,omitempty"`
	ReceiveAge       float64 `json:"rcvage,omitempty"`
	Fragment         uint32  `json:"fragment"`
	SignupFragment   uint32  `json:"signup_fragment"`
	TickPerSecond    uint32  `json:"tps"`
	KeyframeInterval float64 `json:"keyframe_interval,omitempty"`
	TokenRedirect    string  `json:"token_redirect,omitempty"`
	Map              string  `json:"map"`
	Protocol         uint8   `json:"protocol"`
}
