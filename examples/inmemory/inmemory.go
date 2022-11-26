package inmemory

import (
	"sync"
	"time"

	"github.com/FlowingSPDG/gotv-plus-go/gotv"
)

//
// In-memory GOTV+ Engine example
//
// This example does not handle any delay, caching, or hidden option features.
// It only gives you fragment client requested.

var _ gotv.Store = (*InMemory)(nil)
var _ gotv.Broadcaster = (*InMemory)(nil)

// InMemory RAM based GOTV+ Broadcasting Engine
type InMemory struct {
	sync.RWMutex
	password string            // password is Engine-global
	match    map[string]*match // key=token value=match
	delay    int               // frag delay
}

// match SYNC should NOT belong to match
type match struct {
	sync.RWMutex
	ReceiveAge     time.Time
	SignupFragment int
	TickPerSecond  float64
	Protocol       int
	Start          map[int]*gotv.StartFrame // key=fragment_number
	Full           map[int]*gotv.FullFrame  // key=fragment_number
	LastFull       int
	LastDelta      int
	Delta          map[int]*gotv.DeltaFrame // key=fragment_number
	Map            string
}

func (m *InMemory) newMatchIfEmpty(token string) {
	if _, ok := m.match[token]; !ok {
		m.match[token] = &match{
			ReceiveAge:     time.Time{},
			SignupFragment: 0,
			TickPerSecond:  0,
			Protocol:       0,
			Start:          map[int]*gotv.StartFrame{},
			Full:           map[int]*gotv.FullFrame{},
			Delta:          map[int]*gotv.DeltaFrame{},
		}
	}
}

func (m *InMemory) isMatchExist(token string) bool {
	_, ok := m.match[token]
	return ok
}

// Auth implements gotv.Store
func (m *InMemory) Auth(token string, auth string) error {
	if auth != m.password {
		return gotv.ErrInvalidAuth
	}
	return nil
}

func (m *InMemory) isSyncReady(token string, fragment int) bool {
	match, ok := m.match[token]
	if !ok {
		return false
	}
	if _, ok := match.Full[fragment]; !ok {
		return false
	}
	if _, ok := match.Delta[fragment]; !ok {
		return false
	}
	return true
}

// GetSyncLatest implements gotv.Broadcaster
func (m *InMemory) GetSyncLatest(token string) (gotv.Sync, error) {
	m.RLock()
	defer m.RUnlock()
	if !m.isMatchExist(token) {
		return gotv.Sync{}, gotv.ErrMatchNotFound
	}
	match, ok := m.match[token]
	if !ok {
		return gotv.Sync{}, gotv.ErrMatchNotFound
	}
	if !m.isSyncReady(token, match.LastFull-m.delay) {
		return gotv.Sync{}, gotv.ErrFragmentNotFound
	}
	return gotv.Sync{
		Tick:             match.Full[match.LastFull-m.delay].Tick,
		Endtick:          match.Delta[match.LastDelta-m.delay].EndTick,
		RealTimeDelay:    time.Since(match.Full[match.LastFull-m.delay].At).Seconds(),
		ReceiveAge:       time.Since(match.ReceiveAge).Seconds(),
		Fragment:         match.LastFull - m.delay,
		SignupFragment:   match.SignupFragment,
		TickPerSecond:    int(match.TickPerSecond),
		KeyframeInterval: 3, // ?
		TokenRedirect:    "token/" + token,
		Map:              match.Map,
		Protocol:         match.Protocol,
	}, nil
}

// GetSync implements gotv.Broadcaster
func (m *InMemory) GetSync(token string, fragment int) (gotv.Sync, error) {
	m.RLock()
	defer m.RUnlock()
	if !m.isMatchExist(token) {
		return gotv.Sync{}, gotv.ErrMatchNotFound
	}
	match, ok := m.match[token]
	if !ok {
		return gotv.Sync{}, gotv.ErrMatchNotFound
	}
	if !m.isSyncReady(token, fragment) {
		return gotv.Sync{}, gotv.ErrFragmentNotFound
	}
	now := time.Now()
	return gotv.Sync{
		Tick:             match.Full[fragment].Tick,
		Endtick:          match.Delta[fragment].EndTick,
		RealTimeDelay:    (now.Sub(match.Full[fragment].At)).Seconds(),
		ReceiveAge:       (now.Sub(match.ReceiveAge)).Seconds(),
		Fragment:         fragment,
		SignupFragment:   match.SignupFragment,
		TickPerSecond:    int(match.TickPerSecond),
		KeyframeInterval: 3, // ?
		// TokenRedirect:    token,
		Map:      match.Map,
		Protocol: match.Protocol,
	}, nil
}

// GetDelta implements gotv.Broadcaster
func (m *InMemory) GetDelta(token string, fragment int) ([]byte, error) {
	m.RLock()
	defer m.RUnlock()
	if !m.isMatchExist(token) {
		return nil, gotv.ErrFragmentNotFound
	}
	match, ok := m.match[token]
	if !ok {
		return nil, gotv.ErrMatchNotFound
	}
	b, ok := match.Delta[fragment]
	if !ok {
		return nil, gotv.ErrMatchNotFound
	}
	return b.Body, nil
}

// GetFull implements gotv.Broadcaster
func (m *InMemory) GetFull(token string, fragment int) ([]byte, error) {
	m.RLock()
	defer m.RUnlock()
	if !m.isMatchExist(token) {
		return nil, gotv.ErrFragmentNotFound
	}
	match, ok := m.match[token]
	if !ok {
		return nil, gotv.ErrMatchNotFound
	}
	b, ok := match.Full[fragment]
	if !ok {
		return nil, gotv.ErrMatchNotFound
	}
	return b.Body, nil
}

// GetStart implements gotv.Broadcaster
func (m *InMemory) GetStart(token string, fragment int) ([]byte, error) {
	m.RLock()
	defer m.RUnlock()
	if !m.isMatchExist(token) {
		return nil, gotv.ErrFragmentNotFound
	}
	match, ok := m.match[token]
	if !ok {
		return nil, gotv.ErrMatchNotFound
	}
	b, ok := match.Start[fragment]
	if !ok {
		return nil, gotv.ErrMatchNotFound
	}
	return b.Body, nil
}

// OnStart implements gotv.Store
func (m *InMemory) OnStart(token string, fragment int, f gotv.StartFrame) error {
	m.Lock()
	defer m.Unlock()
	m.newMatchIfEmpty(token)
	m.match[token].Start[fragment] = &f
	m.match[token].SignupFragment = fragment
	m.match[token].TickPerSecond = f.Tps
	m.match[token].Protocol = f.Protocol
	m.match[token].Map = f.Map
	return nil
}

// OnFull implements gotv.Store
func (m *InMemory) OnFull(token string, fragment int, f gotv.FullFrame) error {
	m.Lock()
	defer m.Unlock()
	if !m.isMatchExist(token) {
		return gotv.ErrMatchNotFound
	}
	m.match[token].Full[fragment] = &f
	m.match[token].LastFull = fragment
	m.match[token].ReceiveAge = time.Now()
	return nil
}

// OnDelta implements gotv.Store
func (m *InMemory) OnDelta(token string, fragment int, f gotv.DeltaFrame) error {
	m.Lock()
	defer m.Unlock()
	if !m.isMatchExist(token) {
		return gotv.ErrMatchNotFound
	}
	m.match[token].Delta[fragment] = &f
	m.match[token].LastDelta = fragment
	return nil
}

// NewInmemoryGOTV Get new pointer of inMemory GOTV+ Engine
func NewInmemoryGOTV(password string) *InMemory {
	return &InMemory{
		RWMutex:  sync.RWMutex{},
		password: password,
		match:    map[string]*match{},
		delay:    8,
	}
}
