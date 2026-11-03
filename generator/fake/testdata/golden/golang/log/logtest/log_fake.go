package logtest

import (
	"context"
	"sync"

	"example.com/shop/log"
)

// LogFake implements Log in memory, and is safe for concurrent use. It
// appends the values of its methods to a list. The test of LogFake applies
// the suite of Log to it.
type LogFake struct {
	entries []log.Entry
	mu      sync.Mutex
}

// The compiler checks that LogFake implements Log.
var _ log.Log = (*LogFake)(nil)

// NewLogFake returns an empty LogFake.
func NewLogFake() *LogFake {
	return &LogFake{}
}

// Append appends e, and returns the length of the list as its position.
func (f *LogFake) Append(ctx context.Context, e log.Entry) (int64, error) {
	var r0 int64
	if err := ctx.Err(); err != nil {
		return r0, err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.entries = append(f.entries, e)
	return int64(len(f.entries)), nil
}
