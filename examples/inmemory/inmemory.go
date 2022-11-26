package inmemory

import "github.com/FlowingSPDG/gotv-plus-go/gotv"

var _ gotv.Store = (*InMemory)(nil)
var _ gotv.Broadcaster = (*InMemory)(nil)

// InMemory RAM based GOTV+ Broadcasting Engine
type InMemory struct{}

// GetDelta implements gotv.Broadcaster
func (*InMemory) GetDelta(token string, fragment int, endtick int) ([]byte, error) {
	panic("unimplemented")
}

// GetFull implements gotv.Broadcaster
func (*InMemory) GetFull(token string, fragment int, tick int) ([]byte, error) {
	panic("unimplemented")
}

// GetStart implements gotv.Broadcaster
func (*InMemory) GetStart(token string, fragment int) ([]byte, error) {
	panic("unimplemented")
}

// GetSync implements gotv.Broadcaster
func (*InMemory) GetSync(token string) error {
	panic("unimplemented")
}

// Auth implements gotv.Store
func (*InMemory) Auth(token string, auth string) error {
	panic("unimplemented")
}

// OnDelta implements gotv.Store
func (*InMemory) OnDelta(token string, f gotv.DeltaFrame) error {
	panic("unimplemented")
}

// OnFull implements gotv.Store
func (*InMemory) OnFull(token string, f gotv.FullFrame) error {
	panic("unimplemented")
}

// OnStart implements gotv.Store
func (*InMemory) OnStart(token string, f gotv.StartFrame) error {
	panic("unimplemented")
}

// NewInmemoryGOTV Get new pointer of inMemory GOTV+ Engine
func NewInmemoryGOTV() *InMemory {
	return &InMemory{}
}
