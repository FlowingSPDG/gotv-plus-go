package models

import (
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	pb "github.com/FlowingSPDG/gotv-plus-go/server/src/grpc/protogen"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

// Match Match itself.
type Match struct {
	sync.Mutex
	ID    string // Manually tagged ID.
	Token string // Match token
	Auth  string // auth for POST auths

	Delay uint32

	Startframe  map[uint32]*Startframe // start frame data
	Fullframes  map[uint32]*Fullframe  // full frame data
	Deltaframes map[uint32]*Deltaframe // delta frame data

	SignupFragment uint32 // sign up fragment for /sync
	Tps            uint32 // tickrate per secs for /sync
	Map            string // map for /sync
	Protocol       uint8  // protocol for /sync

	// RtDelay uint8  // Real-time delay: delay of this fragment from real-time, in seconds
	// RcVage  uint8  // Receive age: how many seconds since relay last received data from game server
	Latest uint32 // latest fragment number
}

// RegisterStartFrame Register Startframe to match.
func (m *Match) RegisterStartFrame(fragment uint32, start *Startframe, tps uint32) error {
	if m.Startframe == nil {
		m.Startframe = make(map[uint32]*Startframe)
	}
	m.Lock()
	defer m.Unlock()
	m.Startframe[uint32(fragment)] = start
	m.SignupFragment = fragment
	m.Tps = tps
	return nil
}

// RegisterFullFrame Register Fullframe to match.
func (m *Match) RegisterFullFrame(fragment uint32, full *Fullframe) error {
	if m.Fullframes == nil {
		m.Fullframes = make(map[uint32]*Fullframe)
	}
	m.Lock()
	defer m.Unlock()
	m.Latest = fragment
	m.Fullframes[uint32(fragment)] = full
	return nil
}

// RegisterDeltaFrame Register Deltaframe to match.
func (m *Match) RegisterDeltaFrame(fragment uint32, delta *Deltaframe) error {
	if m.Deltaframes == nil {
		m.Deltaframes = make(map[uint32]*Deltaframe)
	}
	m.Lock()
	defer m.Unlock()
	m.Latest = fragment
	m.Deltaframes[uint32(fragment)] = delta
	return nil
}

// GetStartFrame Get START frame by fragnumber.
func (m *Match) GetStartFrame(fragnumber uint32) (*Startframe, error) {
	m.Lock()
	defer m.Unlock()
	if f, ok := m.Startframe[fragnumber]; ok {
		return f, nil
	}
	return nil, fmt.Errorf("Not Found")
}

// GetFullFrame Get FULL frame by fragnumber.
func (m *Match) GetFullFrame(fragnumber uint32) (*Fullframe, error) {
	m.Lock()
	defer m.Unlock()
	if f, ok := m.Fullframes[fragnumber]; ok {
		return f, nil
	}
	return nil, fmt.Errorf("Not Found")
}

// GetDeltaFrame Get DELTA frame by fragnumber.
func (m *Match) GetDeltaFrame(fragnumber uint32) (*Deltaframe, error) {
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
		Tick:           delayed.Tick,
		TokenRedirect:  "token/" + m.Token,
		Endtick:        d.EndTick,
		RealTimeDelay:  time.Since(delayed.At).Seconds(),
		ReceiveAge:     time.Since(latest.At).Seconds(),
		Fragment:       fragnumber - m.Delay,
		SignupFragment: m.SignupFragment,
		TickPerSecond:  m.Tps,
		// KeyframeInterval: 3,
		Map:      m.Map,
		Protocol: m.Protocol,
	}

	return s, nil
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
			Tick:     v.Tick,
			Body:     v.Body,
			At:       t,
		})
	}

	for k, v := range m.Deltaframes {
		binary.DeltaFrame = append(binary.DeltaFrame, &pb.DeltaFrameBinary{
			Fragment: k,
			Endtick:  v.EndTick,
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
