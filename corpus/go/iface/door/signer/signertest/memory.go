// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package signertest

import (
	"context"
	"strings"
	"sync"

	"go.dokimi.dev/conform"

	"go.dokimi.dev/dokimi/corpus/go/iface/door/signer"
)

// signature is the separator between the key and the signature of a request.
const signature = "?signature="

// Memory is the signer in memory. It signs a request with its secret, and its
// backend serves the body of a request whose signature is the secret.
//
//dokimi:subject signer.Signer new=NewMemory options=MemoryOptions
type Memory struct {
	bodies map[string]string
	secret string
	mu     sync.Mutex
}

// NewMemory returns an empty signer in memory.
func NewMemory() *Memory { return &Memory{bodies: map[string]string{}, secret: "secret"} }

// Put stores the body under the key. It returns the cause of ctx when ctx has
// ended.
func (m *Memory) Put(ctx context.Context, key, body string) error {
	if ctx.Err() != nil {
		return context.Cause(ctx)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.bodies[key] = body
	return nil
}

// Sign returns the key with the signature of the secret, and
// signer.ErrNotFound for a key without a body. It returns the cause of ctx
// when ctx has ended.
func (m *Memory) Sign(ctx context.Context, key string) (string, error) {
	return m.signed(ctx, key, m.secret)
}

// Backend returns the backend of the signer, which serves the bodies of its
// signed requests.
func (m *Memory) Backend() Backend { return backend{m: m} }

// signed returns the key with the signature secret, as Sign does.
func (m *Memory) signed(ctx context.Context, key, secret string) (string, error) {
	if ctx.Err() != nil {
		return "", context.Cause(ctx)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, found := m.bodies[key]; !found {
		return "", signer.ErrNotFound
	}
	return key + signature + secret, nil
}

// backend is the backend of a signer in memory.
type backend struct {
	m *Memory
}

// Serve returns the body under the key of the request, and false for a
// request whose signature is not the secret of the signer.
func (b backend) Serve(request string) (string, bool) {
	key, secret, signed := strings.Cut(request, signature)
	b.m.mu.Lock()
	defer b.m.mu.Unlock()
	if !signed || secret != b.m.secret {
		return "", false
	}
	body, found := b.m.bodies[key]
	return body, found
}

// refusing is the backend of a subject that is no signer in memory. It
// refuses every request.
type refusing struct{}

// Serve refuses the request.
func (refusing) Serve(string) (string, bool) { return "", false }

// MemoryOptions returns the door of the backend of a signer in memory.
func MemoryOptions() []SignerOption {
	return []SignerOption{conform.Supplies(conform.Supply(backendOf))}
}

// backendOf returns the backend of the signer in memory s.
func backendOf(_ *SignerHandle, s signer.Signer) Backend {
	if m, isMemory := s.(*Memory); isMemory {
		return m.Backend()
	}
	return refusing{}
}
