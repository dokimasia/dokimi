// Package cells is a cell of one value that a component keeps on the disk of
// its host.
package cells

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"strconv"
	"sync/atomic"

	"go.dokimi.dev/host"
)

// valueFile is the file of the value of a cell.
const valueFile = "value"

// Cell keeps one value in a file of the disk of its host. Each Set writes a
// file of its own and renames it to the file of the value, so concurrent Sets
// need no lock.
type Cell struct {
	env  *host.Env
	sets atomic.Int64
}

// Open returns the cell over the environment env. It syncs the directory of
// the disk first, so a rename that a crash of the process left pending is
// durable before the cell serves a read.
func Open(env *host.Env) *Cell {
	_ = env.Disk().SyncDir(".")
	return &Cell{env: env}
}

// Set makes v the value of the cell, and returns once v is durable.
func (c *Cell) Set(_ context.Context, v string) error {
	disk := c.env.Disk()
	name := valueFile + "." + strconv.FormatInt(c.sets.Add(1), 10)
	f, err := disk.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	_, err = io.WriteString(f, v)
	if err == nil {
		err = f.Sync()
	}
	if err = errors.Join(err, f.Close()); err != nil {
		return err
	}
	if err := disk.Rename(name, valueFile); err != nil {
		return err
	}
	return disk.SyncDir(".")
}

// Get returns the value of the cell, and the empty value for a cell without
// one.
func (c *Cell) Get(_ context.Context) (string, error) {
	f, err := c.env.Disk().OpenFile(valueFile, os.O_RDONLY, 0)
	if errors.Is(err, fs.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	b, err := io.ReadAll(f)
	return string(b), errors.Join(err, f.Close())
}
