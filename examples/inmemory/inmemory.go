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
	Latest         int
	SignupFragment int
	TickPerSecond  float64
	Protocol       int
	Start          map[int]*gotv.StartFrame // key=fragment_number
	Fragments      map[int]*gotv.Fragment   // key=fragment_number
	Map            string
}

func (m *InMemory) newMatchIfEmpty(token string) {
	if _, ok := m.match[token]; !ok {
		m.match[token] = &match{
			RWMutex:        sync.RWMutex{},
			ReceiveAge:     time.Time{},
			SignupFragment: 0,
			TickPerSecond:  0,
			Protocol:       0,
			Start:          map[int]*gotv.StartFrame{},
			Fragments:      map[int]*gotv.Fragment{},
			Map:            "",
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
	if _, ok := match.Fragments[fragment]; !ok {
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
	if !m.isSyncReady(token, match.Latest-m.delay) {
		return gotv.Sync{}, gotv.ErrFragmentNotFound
	}
	return gotv.Sync{
		Tick:             match.Fragments[match.Latest-m.delay].Tick,
		Endtick:          match.Fragments[match.Latest-m.delay].EndTick,
		RealTimeDelay:    time.Since(match.Fragments[match.Latest-m.delay].At).Seconds(),
		ReceiveAge:       time.Since(match.ReceiveAge).Seconds(),
		Fragment:         match.Latest - m.delay,
		SignupFragment:   match.SignupFragment,
		TickPerSecond:    int(match.TickPerSecond),
		KeyframeInterval: 3, // ?
		// TokenRedirect:    "token/" + token,
		Map:      match.Map,
		Protocol: match.Protocol,
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
		Tick:             match.Fragments[fragment].Tick,
		Endtick:          match.Fragments[fragment].EndTick,
		RealTimeDelay:    (now.Sub(match.Fragments[fragment].At)).Seconds(),
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
	b, ok := match.Fragments[fragment]
	if !ok {
		return nil, gotv.ErrMatchNotFound
	}
	return b.Delta, nil
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
	b, ok := match.Fragments[fragment]
	if !ok {
		return nil, gotv.ErrMatchNotFound
	}
	return b.Full, nil
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
func (m *InMemory) OnFull(token string, fragment int, tick int, at time.Time, b []byte) error {
	m.Lock()
	defer m.Unlock()
	if !m.isMatchExist(token) {
		return gotv.ErrMatchNotFound
	}
	if m.match[token].Fragments[fragment] == nil {
		m.match[token].Fragments[fragment] = &gotv.Fragment{}
	}
	m.match[token].Fragments[fragment].At = at
	m.match[token].Fragments[fragment].Tick = tick
	m.match[token].Fragments[fragment].Full = b
	m.match[token].Latest = fragment
	m.match[token].ReceiveAge = time.Now()
	return nil
}

// OnDelta implements gotv.Store
func (m *InMemory) OnDelta(token string, fragment int, endtick int, at time.Time, final bool, b []byte) error {
	m.Lock()
	defer m.Unlock()
	if !m.isMatchExist(token) {
		return gotv.ErrMatchNotFound
	}
	if m.match[token].Fragments[fragment] == nil {
		m.match[token].Fragments[fragment] = &gotv.Fragment{}
	}
	m.match[token].Fragments[fragment].EndTick = endtick
	m.match[token].Fragments[fragment].Final = final
	m.match[token].Fragments[fragment].Delta = b
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
