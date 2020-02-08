package models

import (
	"fmt"
	"time"
)

type Match struct {
	Token       string                  // Match token
	Startframe  map[uint32]*Startframe  // start frame data
	Fullframes  map[uint32]*Fullframe   // full frame data
	Deltaframes map[uint32]*Deltaframes // delta frame data

	SignupFragment uint32  // sign up fragment for /sync
	Tps            float64 // tickrate per secs for /sync
	Map            string  // map for /sync
	Protocol       uint8   // protocol for /sync
	Auth           string  // auth for POST auths

	Tick     uint32 // current tick for /sync
	RtDelay  uint8  // current rtdelay for /sync
	RcVage   uint8  // current rtvage for /sync
	Fragment uint32 // latest fragment
}

func (m *Match) GetBody(ftype string, fragnumber uint32) ([]byte, error) {
	switch ftype {
	case "start":
		if f, ok := m.Startframe[fragnumber]; ok {
			return f.Body, nil
		}
		return nil, fmt.Errorf("Not Found")
	case "full":
		if f, ok := m.Fullframes[fragnumber]; ok {
			return f.Body, nil
		}
		return nil, fmt.Errorf("Not Found")
	case "delta":
		if f, ok := m.Deltaframes[fragnumber]; ok {
			return f.Body, nil
		}
		return nil, fmt.Errorf("Not Found")
	}
	return nil, fmt.Errorf("Unknown ftype")
}

type Matches struct {
	Matches map[string]*Match // string=token
	Auth    string
}

func (m *Matches) Register(ms *Match) {
	if m == nil {
		m = &Matches{}
	}
	if m.Matches == nil {
		m.Matches = make(map[string]*Match)
	}
	m.Matches[ms.Token] = ms
}

func (m *Matches) GetTokens() ([]string, error) { // Gets tokens as slice
	if m == nil {
		return nil, fmt.Errorf("m == nil")
	}
	if m.Matches == nil {
		return nil, fmt.Errorf("m.Matches == nil")
	}
	tokens := make([]string, 0, len(m.Matches))
	for _, v := range m.Matches {
		tokens = append(tokens, v.Token)
	}
	return tokens, nil
}

func (m *Matches) GetMatch(token string) (*Match, error) { // Gets tokens
	if m == nil {
		return nil, fmt.Errorf("m == nil")
	}
	if m.Matches == nil {
		return nil, fmt.Errorf("m.Matches == nil")
	}
	if match, ok := m.Matches[token]; ok {
		return match, nil
	}
	return nil, fmt.Errorf("not found")
}

type Startframe struct {
	At   time.Time
	Body []byte
}

type Fullframe struct {
	At   time.Time
	Tick int
	Body []byte
}

type Deltaframes struct {
	Body []byte
}
