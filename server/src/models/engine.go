package models

import (
	"fmt"
	"runtime"
	"sync"
)

// InitMatchEngine Initializes MatchEngine
func InitMatchEngine(auth string, delay uint32) *MatchesEngine {
	return &MatchesEngine{
		Matches: make(map[string]*Match),
		Auth:    auth, //  tv_broadcast_origin_auth "gopher"
		Delay:   delay,
	}
}

// RelegateMatchesByID Relegates matches matching ID, but not matching token
func (m *MatchesEngine) RelegateMatchesByID(id string, token string) error {
	if m.Matches == nil {
		return fmt.Errorf("m.Matches == nil")
	}
	m.Lock()
	defer m.Unlock()
	for i := range m.Matches {
		if m.Matches[i].ID == id && m.Matches[i].Token != token {
			// Make sure that ID will not match any other request, but can match if need to search for ID and Token
			m.Matches[i].ID = m.Matches[i].ID + "/" + m.Matches[i].Token
		}
	}
	return nil
}

// MatchesEngine Match/Fragment Manager engine.
type MatchesEngine struct {
	sync.Mutex
	Matches map[string]*Match // string=token
	Auth    string
	Delay   uint32
}

// Register Register Match.
func (m *MatchesEngine) Register(id string, token string, signupFragment uint32, tps uint32, mapname string, protocol uint8, auth string) {
	if m.Matches == nil {
		m.Matches = make(map[string]*Match)
	}
	m.Lock()
	defer m.Unlock()
	m.Matches[token] = &Match{
		ID:             id,
		Token:          token,
		Startframe:     make(map[uint32]StartFragment),
		Fullframes:     make(map[uint32]FullFragment),
		Deltaframes:    make(map[uint32]DeltaFragment),
		SignupFragment: signupFragment,
		Tps:            tps,
		Map:            mapname,
		Protocol:       protocol,
		Auth:           auth,
	}
}

// Delete Delete Match and run GC.
func (m *MatchesEngine) Delete(token string) error {
	if m.Matches == nil {
		return fmt.Errorf("m.Matches == nil")
	}
	m.Lock()
	defer m.Unlock()
	m.Matches[token] = nil
	delete(m.Matches, token)
	runtime.GC()
	return nil
}

// GetTokens Gets tokens as slice. THIS IS NOT SORTED.
func (m *MatchesEngine) GetTokens() ([]string, error) {
	if m.Matches == nil {
		return nil, fmt.Errorf("m.Matches == nil")
	}
	m.Lock()
	defer m.Unlock()
	tokens := make([]string, 0, len(m.Matches))
	for _, v := range m.Matches {
		tokens = append(tokens, v.Token)
	}
	return tokens, nil
}

// GetAll Gets tokens as slice
func (m *MatchesEngine) GetAll() ([]*Match, error) {
	if m.Matches == nil {
		return nil, fmt.Errorf("m.Matches == nil")
	}
	m.Lock()
	defer m.Unlock()
	matches := make([]*Match, 0, len(m.Matches))
	for _, v := range m.Matches {
		matches = append(matches, v)
	}
	return matches, nil
}

// GetMatchByToken Get Match pointer by token string.
func (m *MatchesEngine) GetMatchByToken(token string) (*Match, error) {
	if m.Matches == nil {
		return nil, fmt.Errorf("m.Matches == nil")
	}
	m.Lock()
	defer m.Unlock()
	if match, ok := m.Matches[token]; ok {
		return match, nil
	}
	return nil, fmt.Errorf("not found")
}

// GetMatchByID Get Match pointer by id string.
func (m *MatchesEngine) GetMatchByID(id string) (*Match, error) { // Gets tokens
	if m.Matches == nil {
		return nil, fmt.Errorf("m.Matches == nil")
	}
	m.Lock()
	defer m.Unlock()
	for _, v := range m.Matches {
		if v.ID == id {
			return v, nil
		}
	}
	return nil, fmt.Errorf("not found")
}
