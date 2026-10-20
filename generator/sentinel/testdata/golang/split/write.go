package split

import "errors"

// ErrWrite reports a failed write.
var ErrWrite = errors.New("split: write failed")
