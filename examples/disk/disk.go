package disk

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/FlowingSPDG/gotv-plus-go/gotv"
)

//
// Disk based GOTV+ Engine example
//
// This example does not handle any delay, caching, or hidden option features.
// It only gives you fragment client requested.

var _ gotv.Store = (*Disk)(nil)
var _ gotv.Broadcaster = (*Disk)(nil)

// Disk fragment disk file based GOTV+ Broadcasting Engine
type Disk struct {
	password string // password is Engine-global
	dir      string // Work dir
}

func (d *Disk) deltaFramePath(token string, fragment int) string {
	return path.Join(d.dir, fmt.Sprintf("%s_%d_delta.bin", token, fragment))
}
func (d *Disk) startFramePath(token string, fragment int) string {
	return path.Join(d.dir, fmt.Sprintf("%s_%d_start.bin", token, fragment))
}
func (d *Disk) fullFramePath(token string, fragment int) string {
	return path.Join(d.dir, fmt.Sprintf("%s_%d_full.bin", token, fragment))
}
func (d *Disk) syncPath(token string) string {
	return path.Join(d.dir, fmt.Sprintf("%s_sync.json", token))
}

// GetDelta implements gotv.Broadcaster
func (d *Disk) GetDelta(token string, fragment int) ([]byte, error) {
	b, err := os.ReadFile(d.deltaFramePath(token, fragment))
	if err != nil {
		if err == os.ErrNotExist {
			return nil, gotv.ErrFragmentNotFound
		}
		return nil, err
	}
	return b, err
}

// GetFull implements gotv.Broadcaster
func (d *Disk) GetFull(token string, fragment int) ([]byte, error) {
	b, err := os.ReadFile(d.fullFramePath(token, fragment))
	if err != nil {
		if err == os.ErrNotExist {
			return nil, gotv.ErrFragmentNotFound
		}
		return nil, err
	}
	return b, err
}

// GetStart implements gotv.Broadcaster
func (d *Disk) GetStart(token string, fragment int) ([]byte, error) {
	b, err := os.ReadFile(d.startFramePath(token, fragment))
	if err != nil {
		if err == os.ErrNotExist {
			return nil, gotv.ErrFragmentNotFound
		}
		return nil, err
	}
	return b, err
}

// GetSync implements gotv.Broadcaster
func (d *Disk) GetSync(token string) (gotv.Sync, error) {
	ret := gotv.Sync{}
	b, err := os.ReadFile(d.syncPath(token))
	if err != nil {
		if err == os.ErrNotExist {
			return ret, gotv.ErrFragmentNotFound
		}
		return ret, err
	}

	if err := json.Unmarshal(b, &ret); err != nil {
		return ret, err
	}
	return ret, nil
}

// OnDelta implements gotv.Store
func (d *Disk) OnDelta(token string, fragment int, df gotv.DeltaFrame) error {
	return os.WriteFile(d.deltaFramePath(token, fragment), df.Body, 0755)
}

// OnFull implements gotv.Store
func (d *Disk) OnFull(token string, fragment int, ff gotv.FullFrame) error {
	return os.WriteFile(d.fullFramePath(token, fragment), ff.Body, 0755)
}

// OnStart implements gotv.Store
func (d *Disk) OnStart(token string, fragment int, sf gotv.StartFrame) error {
	return os.WriteFile(d.startFramePath(token, fragment), sf.Body, 0755)
}

// Auth implements gotv.Store
func (d *Disk) Auth(token string, auth string) error {
	if auth != d.password {
		return gotv.ErrInvalidAuth
	}
	return nil
}

// NewDiskGOTV Get new pointer of Disk GOTV+ Engine
func NewDiskGOTV(password string, dir string) *Disk {
	p := path.Join(".", dir)
	os.MkdirAll(p, 0755)
	return &Disk{
		password: password,
		dir:      p,
	}
}
