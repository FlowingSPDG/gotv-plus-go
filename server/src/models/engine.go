package models

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"

	pb "github.com/FlowingSPDG/gotv-plus-go/server/src/grpc/protogen"
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
func (m *MatchesEngine) Register(ms *Match) {
	if m.Matches == nil {
		m.Matches = make(map[string]*Match)
	}
	m.Lock()
	defer m.Unlock()
	m.Matches[ms.Token] = ms
}

// LoadMatchFromFile Load match from gzip compressed file
func (m *MatchesEngine) LoadMatchFromFile(path string) (string, error) {
	if m.Matches == nil {
		m.Matches = make(map[string]*Match)
	}

	file, err := os.Open(fmt.Sprintf("matches/%s.gz", path))
	if err != nil {
		return "", err
	}
	gzipreader, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	buf := &pb.MatchBinary{}
	bytes, err := ioutil.ReadAll(gzipreader)
	if err != nil {
		return "", err
	}
	err = proto.Unmarshal(bytes, buf)

	match := &Match{
		ID:             buf.Id,
		Token:          buf.Token,
		Startframe:     make(map[uint32]*Startframe),
		Fullframes:     make(map[uint32]*Fullframe),
		Deltaframes:    make(map[uint32]*Deltaframes),
		Tps:            buf.StartFrame[0].Tps,
		Map:            buf.StartFrame[0].Map,
		Protocol:       uint8(buf.StartFrame[0].Protocol),
		Auth:           "", // TODO
		SignupFragment: buf.StartFrame[0].Fragment,
		// Latest:       buf.FullFrame[0].Fragment, // TODO?
	}

	for _, v := range buf.StartFrame {
		// t, _ := ptypes.Timestamp(v.At)
		match.Startframe[v.Fragment] = &Startframe{
			// At:   t,
			At:   time.Now(),
			Body: v.Body,
		}
	}

	fulls := make([]uint32, 0, len(match.Fullframes))
	for _, v := range buf.FullFrame {
		// t, _ := ptypes.Timestamp(v.At)
		match.Fullframes[v.Fragment] = &Fullframe{
			//At:   t,
			At:   time.Now(),
			Tick: v.Tick,
			Body: v.Body,
		}
		fulls = append(fulls, v.Fragment)
	}

	deltas := make([]uint32, 0, len(match.Deltaframes))
	for _, v := range buf.DeltaFrame {
		match.Deltaframes[v.Fragment] = &Deltaframes{
			EndTick: v.Endtick,
			Body:    v.Body,
		}
		deltas = append(deltas, v.Fragment)
	}

	m.Register(match)

	log.Printf("Loaded match from %s. Available Full list : [%v], Delta : [%v]\n", file.Name(), fulls, deltas)
	return match.ID, nil
}

// Delete Delete Match and run GC.
func (m *MatchesEngine) Delete(ms *Match) error {
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
