package gcs

import (
	"time"

	"cloud.google.com/go/storage"
	"github.com/FlowingSPDG/gotv-plus-go/gotv"
)

//
// Google Cloud Storage GOTV+ Engine example
//

var _ gotv.Store = (*CloudStorage)(nil)
var _ gotv.Broadcaster = (*CloudStorage)(nil)

// CloudStorage GCS based GOTV+ Broadcasting Engine
type CloudStorage struct {
	s        *storage.Client // Firebase Storage and Google Cloud Storage is identical
	password string          // password
	delay    int             // frag delay
}

// Auth implements gotv.Store
func (c *CloudStorage) Auth(token string, auth string) error {
	if auth != c.password {
		return gotv.ErrInvalidAuth
	}
	return nil
}

// OnDelta implements gotv.Store
func (c *CloudStorage) OnDelta(token string, fragment int, endtick int, at time.Time, final bool, b []byte) error {
	panic("unimplemented")
}

// OnFull implements gotv.Store
func (c *CloudStorage) OnFull(token string, fragment int, tick int, at time.Time, b []byte) error {
	panic("unimplemented")
}

// OnStart implements gotv.Store
func (c *CloudStorage) OnStart(token string, fragment int, f gotv.StartFrame) error {
	panic("unimplemented")
}

// GetDelta implements gotv.Broadcaster
func (c *CloudStorage) GetDelta(token string, fragment int) ([]byte, error) {
	panic("unimplemented")
}

// GetFull implements gotv.Broadcaster
func (c *CloudStorage) GetFull(token string, fragment int) ([]byte, error) {
	panic("unimplemented")
}

// GetStart implements gotv.Broadcaster
func (c *CloudStorage) GetStart(token string, fragment int) ([]byte, error) {
	panic("unimplemented")
}

// GetSync implements gotv.Broadcaster
func (c *CloudStorage) GetSync(token string, fragment int) (gotv.Sync, error) {
	panic("unimplemented")
}

// GetSyncLatest implements gotv.Broadcaster
func (c *CloudStorage) GetSyncLatest(token string) (gotv.Sync, error) {
	panic("unimplemented")
}

// NewCloudStorageGOTV Get new pointer of GCS GOTV+ Engine
func NewCloudStorageGOTV(s *storage.Client, password string, delay int) *CloudStorage {
	return &CloudStorage{
		s:        s,
		password: password,
		delay:    delay,
	}
}
