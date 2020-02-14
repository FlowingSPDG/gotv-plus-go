package handlers

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
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

func (m *Match) SaveMatchToFile(path string) error {
	log.Printf("Saving match %s to file...\n", m.Token)

	log.Println("Writing start fragmnets...")
	startfile, err := os.Create(fmt.Sprintf("%s/%s_start.tar.gz", path, m.Token))
	if err != nil {
		return err
	}
	defer startfile.Close()

	StartgzipReader := gzip.NewWriter(startfile)
	defer StartgzipReader.Close()

	starttw := tar.NewWriter(StartgzipReader)
	defer starttw.Close()

	for k, v := range m.Startframe {
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

	log.Println("Writing full fragmnets...")
	fullfile, err := os.Create(fmt.Sprintf("%s/%s_full.tar.gz", path, m.Token))
	if err != nil {
		return err
	}
	defer fullfile.Close()
	fullgzipReader := gzip.NewWriter(fullfile)
	defer fullgzipReader.Close()
	fulltw := tar.NewWriter(fullgzipReader)
	defer fulltw.Close()

	// Write fragments
	for k, v := range m.Fullframes {
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

	log.Println("Writing delta fragmnets...")
	deltafile, err := os.Create(fmt.Sprintf("%s/%s_delta.tar.gz", path, m.Token))
	if err != nil {
		return err
	}
	defer deltafile.Close()
	deltagzipReader := gzip.NewWriter(deltafile)
	defer deltagzipReader.Close()
	deltatw := tar.NewWriter(deltagzipReader)
	defer deltatw.Close()

	// Write fragments
	for k, v := range m.Deltaframes {
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
	match := &Match{
		Token:       path,
		Startframe:  make(map[uint32]*Startframe),
		Fullframes:  make(map[uint32]*Fullframe),
		Deltaframes: make(map[uint32]*Deltaframes),
		Tps:         32,         // TODO...
		Map:         "de_dust2", // TODO...
		Protocol:    uint8(4),   // TODO...
		Auth:        "",         // TODO...
		Tick:        uint32(0),  // TODO...
		Fragment:    9999,
	}

	startpath := fmt.Sprintf("matches/%s_start.tar.gz", path)
	fullpath := fmt.Sprintf("matches/%s_full.tar.gz", path)
	deltapath := fmt.Sprintf("matches/%s_delta.tar.gz", path)

	startfile, err := os.Open(startpath)
	if err != nil {
		return "", err
	}
	defer startfile.Close()
	startgzipreader, err := gzip.NewReader(startfile)
	if err != nil {
		return "", err
	}
	defer startgzipreader.Close()
	starttarballreader := tar.NewReader(startgzipreader)
	for {
		Header, err := starttarballreader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		fragment, err := strconv.Atoi(Header.Name)
		if err != nil {
			return "", err
		}
		bin, err := ioutil.ReadAll(starttarballreader)
		if err != nil {
			return "", err
		}
		match.RegisterStartFrame(uint32(fragment), &Startframe{
			At:   time.Now(),
			Body: bin,
		})
		log.Printf("READ START FRAGMENT %s in file %s. %d bytes\n", Header.Name, startpath, Header.Size)
	}

	fullfile, err := os.Open(fullpath)
	if err != nil {
		return "", err
	}
	defer fullfile.Close()
	fullgzipreader, err := gzip.NewReader(fullfile)
	if err != nil {
		return "", err
	}
	defer fullgzipreader.Close()
	fulltarballreader := tar.NewReader(fullgzipreader)
	for {
		Header, err := fulltarballreader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		fragment, err := strconv.Atoi(Header.Name)
		if err != nil {
			return "", err
		}
		bin, err := ioutil.ReadAll(fulltarballreader)
		if err != nil {
			return "", err
		}
		match.RegisterFullFrame(uint32(fragment), &Fullframe{
			At:   time.Now(),
			Body: bin,
			Tick: 1,
		})
		if fragment < int(match.Fragment) {
			err = match.UpdateFragment(uint32(fragment))
			if err != nil {
				return "", err
			}
		}
		log.Printf("READ FULL FRAGMENT %s in file %s. %d bytes\n", Header.Name, fullpath, Header.Size)
	}

	deltafile, err := os.Open(deltapath)
	if err != nil {
		return "", err
	}
	defer deltafile.Close()
	deltagzipreader, err := gzip.NewReader(deltafile)
	if err != nil {
		return "", err
	}
	defer deltagzipreader.Close()
	deltatarballreader := tar.NewReader(deltagzipreader)
	for {
		Header, err := deltatarballreader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		fragment, err := strconv.Atoi(Header.Name)
		if err != nil {
			return "", err
		}
		bin, err := ioutil.ReadAll(deltatarballreader)
		if err != nil {
			return "", err
		}
		match.RegisterDeltaFrame(uint32(fragment), &Deltaframes{
			Body:    bin,
			EndTick: 0,
		})
		log.Printf("READ DELTA FRAGMENT %s in file %s. %d bytes\n", Header.Name, deltapath, Header.Size)
	}

	Matches.Register(match)

	return path, nil
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
