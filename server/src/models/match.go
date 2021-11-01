package models

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Match Match itself.
type Match struct {
	sync.Mutex
	ID    string // Manually tagged ID.
	Token string // Match token
	Auth  string // auth for POST auths

	Delay uint32

	Startframe  map[uint32]StartFragment // start frame data
	Fullframes  map[uint32]FullFragment  // full frame data
	Deltaframes map[uint32]DeltaFragment // delta frame data

	SignupFragment uint32 // sign up fragment for /sync
	Tps            uint32 // tickrate per secs for /sync
	Map            string // map for /sync
	Protocol       uint8  // protocol for /sync

	// RtDelay uint8  // Real-time delay: delay of this fragment from real-time, in seconds
	// RcVage  uint8  // Receive age: how many seconds since relay last received data from game server
	Latest uint32 // latest fragment number
}

// RegisterStartFrame Register Startframe to match.
func (m *Match) RegisterStartFrame(fragment uint32, start StartFragment, tps uint32) error {
	if m.Startframe == nil {
		m.Startframe = make(map[uint32]StartFragment)
	}
	m.Lock()
	defer m.Unlock()
	m.Startframe[uint32(fragment)] = start
	m.SignupFragment = fragment
	m.Tps = tps
	return nil
}

// RegisterFullFrame Register Fullframe to match.
func (m *Match) RegisterFullFrame(fragment uint32, full FullFragment) error {
	if m.Fullframes == nil {
		m.Fullframes = make(map[uint32]FullFragment)
	}
	m.Lock()
	defer m.Unlock()
	m.Latest = fragment
	m.Fullframes[uint32(fragment)] = full
	return nil
}

// RegisterDeltaFrame Register Deltaframe to match.
func (m *Match) RegisterDeltaFrame(fragment uint32, delta DeltaFragment) error {
	if m.Deltaframes == nil {
		m.Deltaframes = make(map[uint32]DeltaFragment)
	}
	m.Lock()
	defer m.Unlock()
	m.Latest = fragment
	m.Deltaframes[uint32(fragment)] = delta
	return nil
}

// GetStartFrame Get START frame by fragnumber.
func (m *Match) GetStartFrame(fragnumber uint32) (StartFragment, error) {
	m.Lock()
	defer m.Unlock()
	if f, ok := m.Startframe[fragnumber]; ok {
		return f, nil
	}
	return nil, fmt.Errorf("Not Found")
}

// GetFullFrame Get FULL frame by fragnumber.
func (m *Match) GetFullFrame(fragnumber uint32) (FullFragment, error) {
	m.Lock()
	defer m.Unlock()
	if f, ok := m.Fullframes[fragnumber]; ok {
		return f, nil
	}
	return nil, fmt.Errorf("Not Found")
}

// GetDeltaFrame Get DELTA frame by fragnumber.
func (m *Match) GetDeltaFrame(fragnumber uint32) (DeltaFragment, error) {
	m.Lock()
	defer m.Unlock()
	if f, ok := m.Deltaframes[fragnumber]; ok {
		return f, nil
	}
	return nil, fmt.Errorf("Not Found")
}

// TagID Tag match ID.
func (m *Match) TagID(id string) {
	m.ID = id
}

// IsSyncReady Check requested fragnumber is ready for /sync request.
func (m *Match) IsSyncReady(fragnumber uint32) bool {
	_, err := m.GetDeltaFrame(fragnumber)
	if err != nil {
		log.Printf("SYNC NOT READY : fragment[%d]\n", fragnumber)
		return false
	}
	_, err = m.GetFullFrame(fragnumber)
	if err != nil {
		log.Printf("SYNC NOT READY : fragment[%d]\n", fragnumber)
		return false
	}
	log.Printf("SYNC READY : fragment[%d]\n", fragnumber)
	return true
}

// Sync Get SyncJSON.
func (m *Match) Sync(fragnumber uint32) (*SyncJSON, error) {
	for {
		if fragnumber > m.Latest {
			return nil, fmt.Errorf("ERROR Fragment not found")
		}
		if m.IsSyncReady(fragnumber) {
			break
		}
		fragnumber--
	}

	delayed, err := m.GetFullFrame(fragnumber - m.Delay)
	if err != nil {
		return nil, err
	}
	log.Printf("FULL TICK[%d]\n", delayed.Tick)

	d, err := m.GetDeltaFrame(fragnumber - m.Delay)
	if err != nil {
		return nil, err
	}
	log.Printf("DELTA TICK[%d]\n", d.EndTick)

	latest, _ := m.GetFullFrame(fragnumber)

	s := &SyncJSON{
		Tick:           delayed.Tick(),
		TokenRedirect:  "token/" + m.Token,
		Endtick:        d.EndTick(),
		RealTimeDelay:  time.Since(delayed.At()).Seconds(),
		ReceiveAge:     time.Since(latest.At()).Seconds(),
		Fragment:       fragnumber - m.Delay,
		SignupFragment: m.SignupFragment,
		TickPerSecond:  m.Tps,
		// KeyframeInterval: 3,
		Map:      m.Map,
		Protocol: m.Protocol,
	}

	return s, nil
}
