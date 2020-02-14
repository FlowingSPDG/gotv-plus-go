package handlers

import (
	"archive/tar"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

var (
	Matches *MatchesEngine
	Auth    string
	Delay   uint32
)

// InitMatchEngine Initializes MatchEngine
func InitMatchEngine(auth string, delay uint32) {
	Matches = &MatchesEngine{
		Matches: make(map[string]*Match),
		Auth:    auth, //  tv_broadcast_origin_auth "gopher"
		Delay:   delay,
	}
}

type Match struct {
	sync.Mutex
	ID          string                  // Manually tagged ID.
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
	if m == nil {
		return nil, fmt.Errorf("Match not found")
	}
	m.Lock()
	defer m.Unlock()
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

func (m *Match) RegisterStartFrame(fragment uint32, start *Startframe) error {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.Startframe[uint32(fragment)]; !ok {
		return fmt.Errorf("Not Found")
	}
	m.Startframe[uint32(fragment)] = start
	return nil
}

func (m *Match) RegisterFullFrame(fragment uint32, full *Fullframe) error {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.Fullframes[uint32(fragment)]; !ok {
		return fmt.Errorf("Not Found")
	}
	m.Fullframes[uint32(fragment)] = full
	return nil
}

func (m *Match) RegisterDelstaFrame(fragment uint32, delta *Deltaframes) error {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.Deltaframes[uint32(fragment)]; !ok {
		return fmt.Errorf("Not Found")
	}
	m.Deltaframes[uint32(fragment)] = delta
	return nil
}

func (m *Match) GetFullFrame(fragnumber uint32) (*Fullframe, error) {
	if m == nil {
		return nil, fmt.Errorf("Match not found")
	}
	m.Lock()
	defer m.Unlock()
	if f, ok := m.Fullframes[fragnumber]; ok {
		return f, nil
	}
	return nil, fmt.Errorf("Not Found")
}

func (m *Match) TagID(id string) error {
	if m == nil {
		return fmt.Errorf("Match not found")
	}
	m.ID = id
	return nil
}

func (m *Match) SaveMatchToFile(path string) error {
	if m == nil {
		return fmt.Errorf("Match not found")
	}
	log.Printf("Saving match %s to file...\n", m.Token)

	log.Printf("Writing start fragmnets... %v\n", m.Startframe)
	startfile, err := os.Create(fmt.Sprintf("%s/%s_start.tar", path, m.Token))
	if err != nil {
		return err
	}
	defer startfile.Close()
	starttw := tar.NewWriter(startfile)
	defer starttw.Close()

	for k, v := range m.Startframe {
		log.Printf("Writing start fragment : [%v]\n", v.Body)
		// Header info
		StartHeader := &tar.Header{
			Name: fmt.Sprintf("%d", k),
			Size: int64(len(v.Body)),
		}
		// Write header
		if err := starttw.WriteHeader(StartHeader); err != nil {
			return err
		}
		// Write fragments
		if _, err := starttw.Write(v.Body); err != nil {
			return err
		}
	}

	log.Printf("Writing full fragmnets... %v\n", m.Fullframes)
	fullfile, err := os.Create(fmt.Sprintf("%s/%s_full.tar", path, m.Token))
	if err != nil {
		return err
	}
	defer fullfile.Close()
	fulltw := tar.NewWriter(fullfile)
	defer fulltw.Close()

	// Write fragments
	for k, v := range m.Fullframes {
		log.Printf("Writing full fragment : [%v]\n", v.Body)
		// Header info
		FullHeader := &tar.Header{
			Name: fmt.Sprintf("%d", k),
			Size: int64(len(v.Body)),
		}
		// Write header
		if err := fulltw.WriteHeader(FullHeader); err != nil {
			return err
		}
		if _, err := fulltw.Write(v.Body); err != nil {
			return err
		}
	}

	log.Printf("Writing delta fragmnets... %v\n", m.Deltaframes)
	deltafile, err := os.Create(fmt.Sprintf("%s/%s_delta.tar", path, m.Token))
	if err != nil {
		return err
	}
	defer deltafile.Close()
	deltatw := tar.NewWriter(deltafile)
	defer deltatw.Close()

	// Write fragments
	for k, v := range m.Deltaframes {
		log.Printf("Writing delta fragment : [%v]\n", v.Body)
		// Header info
		DeltaHeader := &tar.Header{
			Name: fmt.Sprintf("%d", k),
			Size: int64(len(v.Body)),
		}
		// Write header
		if err := deltatw.WriteHeader(DeltaHeader); err != nil {
			return err
		}
		if _, err := deltatw.Write(v.Body); err != nil {
			return err
		}
	}
	return nil
}

type MatchesEngine struct {
	sync.Mutex
	Matches map[string]*Match // string=token
	Auth    string
	Delay   uint32
}

func (m *MatchesEngine) Register(ms *Match) {
	if m == nil {
		m = &MatchesEngine{}
	}
	if m.Matches == nil {
		m.Matches = make(map[string]*Match)
	}
	m.Lock()
	defer m.Unlock()
	m.Matches[ms.Token] = ms
}

func (m *MatchesEngine) Delete(ms *Match) error {
	if m == nil {
		return fmt.Errorf("m == nil")
	}
	if m.Matches == nil {
		return fmt.Errorf("m.Matches == nil")
	}
	m.Lock()
	defer m.Unlock()
	m.Matches[ms.Token] = nil
	delete(m.Matches, ms.Token)
	runtime.GC()
	return nil
}

func (m *MatchesEngine) GetTokens() ([]string, error) { // Gets tokens as slice
	if m == nil {
		return nil, fmt.Errorf("m == nil")
	}
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

func (m *MatchesEngine) GetAll() ([]*Match, error) { // Gets tokens as slice
	if m == nil {
		return nil, fmt.Errorf("m == nil")
	}
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

func (m *MatchesEngine) GetMatchByToken(token string) (*Match, error) { // Gets tokens
	if m == nil {
		return nil, fmt.Errorf("m == nil")
	}
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

func (m *MatchesEngine) GetMatchByID(id string) (*Match, error) { // Gets tokens
	if m == nil {
		return nil, fmt.Errorf("m == nil")
	}
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
	Body    []byte
	EndTick int // ??
}
