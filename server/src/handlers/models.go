package handlers

import (
	"compress/gzip"
	"fmt"
	pb "github.com/FlowingSPDG/gotv-plus-go/server/src/grpc/protogen"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"io/ioutil"
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

func (m *Match) UpdateFragment(fragnumber uint32) error {
	m.Lock()
	m.Fragment = fragnumber
	m.Unlock()
	return nil
}

func (m *Match) GetBody(ftype string, fragnumber uint32) ([]byte, error) {
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
	if m.Startframe == nil {
		m.Startframe = make(map[uint32]*Startframe)
	}
	m.Lock()
	defer m.Unlock()
	m.Startframe[uint32(fragment)] = start
	return nil
}

func (m *Match) RegisterFullFrame(fragment uint32, full *Fullframe) error {
	if m.Fullframes == nil {
		m.Fullframes = make(map[uint32]*Fullframe)
	}
	m.Lock()
	defer m.Unlock()
	m.Fullframes[uint32(fragment)] = full
	return nil
}

func (m *Match) RegisterDeltaFrame(fragment uint32, delta *Deltaframes) error {
	if m.Deltaframes == nil {
		m.Deltaframes = make(map[uint32]*Deltaframes)
	}
	m.Lock()
	defer m.Unlock()
	m.Deltaframes[uint32(fragment)] = delta
	return nil
}

func (m *Match) GetFullFrame(fragnumber uint32) (*Fullframe, error) {
	m.Lock()
	defer m.Unlock()
	if f, ok := m.Fullframes[fragnumber]; ok {
		return f, nil
	}
	return nil, fmt.Errorf("Not Found")
}

func (m *Match) TagID(id string) error {
	m.ID = id
	return nil
}

func (m *Match) SaveMatchToFile(filename string) error {
	log.Printf("Saving match %s to file...\n", m.Token)

	file, err := os.Create(fmt.Sprintf("matches/%s.gz", filename))
	if err != nil {
		return err
	}
	defer file.Close()

	binary := &pb.MatchBinary{
		Id:             m.ID,
		Token:          m.Token,
		SignupFragment: m.SignupFragment,
		StartFrame:     make([]*pb.StartFrameBinary, 0, len(m.Startframe)),
		FullFrame:      make([]*pb.FullFrameBinary, 0, len(m.Fullframes)),
		DeltaFrame:     make([]*pb.DeltaFrameBinary, 0, len(m.Deltaframes)),
	}
	for k, v := range m.Startframe {
		t, _ := ptypes.TimestampProto(v.At)
		binary.StartFrame = append(binary.StartFrame, &pb.StartFrameBinary{
			Fragment: k,
			Tick:     m.Tick,
			Tps:      m.Tps,
			Map:      m.Map,
			Protocol: uint32(m.Protocol),
			Body:     v.Body,
			At:       t,
		})
	}

	for k, v := range m.Fullframes {
		t, _ := ptypes.TimestampProto(v.At)
		binary.FullFrame = append(binary.FullFrame, &pb.FullFrameBinary{
			Fragment: k,
			Tick:     m.Tick,
			Body:     v.Body,
			At:       t,
		})
	}

	for k, v := range m.Deltaframes {
		binary.DeltaFrame = append(binary.DeltaFrame, &pb.DeltaFrameBinary{
			Fragment: k,
			Endtick:  uint32(v.EndTick),
			Body:     v.Body,
		})
	}

	data, err := proto.Marshal(binary)
	if err != nil {
		return err
	}
	gzipwriter := gzip.NewWriter(file)
	defer gzipwriter.Close()
	totalbytes, err := gzipwriter.Write(data)
	if err != nil {
		return err
	}
	log.Printf("Writed to %s. %dbytes\n", file.Name(), totalbytes)
	return nil

}

type MatchesEngine struct {
	sync.Mutex
	Matches map[string]*Match // string=token
	Auth    string
	Delay   uint32
}

func (m *MatchesEngine) Register(ms *Match) {
	if m.Matches == nil {
		m.Matches = make(map[string]*Match)
	}
	m.Lock()
	defer m.Unlock()
	m.Matches[ms.Token] = ms
}

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
		Tick:           buf.StartFrame[0].Tick,
		SignupFragment: buf.StartFrame[0].Fragment,
		Fragment:       buf.FullFrame[0].Fragment, // TODO?
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
			Tick: int(v.Tick),
			Body: v.Body,
		}
		fulls = append(fulls, v.Fragment)
	}

	deltas := make([]uint32, 0, len(match.Deltaframes))
	for _, v := range buf.DeltaFrame {
		match.Deltaframes[v.Fragment] = &Deltaframes{
			EndTick: int(v.Endtick),
			Body:    v.Body,
		}
		deltas = append(deltas, v.Fragment)
	}

	Matches.Register(match)

	log.Printf("Loaded match from %s. Available Full list : [%v], Delta : [%v]\n", file.Name(), fulls, deltas)
	return match.ID, nil
}

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

func (m *MatchesEngine) GetTokens() ([]string, error) { // Gets tokens as slice
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
