// Package buffered reads records into the buffers of their requests, and
// returns the buffer of a request with every error.
package buffered

import "context"

// Request is a read of the record of a key into a buffer.
type Request struct {
	Key string
	Dst []byte
}

// Result is the record of a read and its size.
type Result struct {
	Record []byte
	Size   int
}

// Records reads the records of keys.
//
//dokimi:suite
type Records interface {
	// Read appends the record of req.Key to req.Dst. It returns req.Dst as
	// the record of its result with every error.
	//
	//dokimi:mixin returns-on-error field=Dst result=Record
	Read(ctx context.Context, req *Request) (Result, error)

	// Fill appends the record of req.Key to req.Dst, and returns the buffer.
	// It returns req.Dst unchanged with every error.
	//
	//dokimi:mixin returns-on-error field=Dst
	Fill(ctx context.Context, req Request) ([]byte, error)

	// Peek returns the size of the record of req.Key. The field of its mixin
	// is a field of its result, and no field of its request.
	//
	//dokimi:mixin returns-on-error field=Size
	Peek(ctx context.Context, req Request) (Result, error)
}
