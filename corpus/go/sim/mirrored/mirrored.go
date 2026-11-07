// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package mirrored is a register of one value whose primary keeps the value
// on the disk of its host and mirrors it to a backup over a stream. The
// clients call the primary. The primary sends its value to the backup again
// every Period, so the backup has the value of the primary once the faults
// end. Its nodes are Go components that run on the hosts of a simulation.
package mirrored

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"go.dokimi.dev/host"
)

// Port is the port on which the backup takes the values of the primary.
const Port = 7000

// Period is the time between two sends of the value of the primary to the
// backup, on the monotonic clock of the host of the primary.
const Period = 100 * time.Millisecond

// valueFile is the file of the value of a node.
const valueFile = "value"

// ack is the line with which the backup acknowledges a value once the value
// is durable on its disk.
const ack = "ok\n"

// Node is a node of the register: the primary, which mirrors each value to
// the host of the backup, or the backup.
type Node struct {
	env    *host.Env
	backup string
	writes atomic.Int64
}

// Primary returns the open of the primary that mirrors to the host named
// backup. The open makes the last rename of a value durable, which a crash of
// the process can leave pending, and starts the sends of the value.
func Primary(backup string) func(env *host.Env) *Node {
	return func(env *host.Env) *Node {
		n := &Node{env: env, backup: backup}
		_ = env.Disk().SyncDir(".")
		env.Go(n.resend)
		return n
	}
}

// Seeded returns the open of Primary over a disk that starts with value.
func Seeded(value, backup string) func(env *host.Env) *Node {
	open := Primary(backup)
	return func(env *host.Env) *Node {
		n := open(env)
		_ = n.store(value)
		return n
	}
}

// Backup opens the backup over the environment env. It makes the last rename
// of a value durable, and takes the values of the primary on Port.
func Backup(env *host.Env) *Node {
	n := &Node{env: env}
	_ = env.Disk().SyncDir(".")
	if l, err := env.Listen(Port); err == nil {
		env.Go(func() { n.accept(l) })
	}
	return n
}

// Set mirrors v to the backup, and then makes v the value of the primary. It
// returns once v is durable on both disks. A Set that cannot mirror v leaves
// the value of the primary as it was.
func (n *Node) Set(ctx context.Context, v string) error {
	if err := n.mirror(ctx, v); err != nil {
		return err
	}
	return n.store(v)
}

// Get returns the value of the node, and the empty value for a node without
// one.
func (n *Node) Get(context.Context) (string, error) {
	f, err := n.env.Disk().OpenFile(valueFile, os.O_RDONLY, 0)
	if errors.Is(err, fs.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("mirrored: open the value: %w", err)
	}
	b, err := io.ReadAll(f)
	if err = errors.Join(err, f.Close()); err != nil {
		return "", fmt.Errorf("mirrored: read the value: %w", err)
	}
	return string(b), nil
}

// mirror sends v to the backup over a stream of its own, and returns once the
// backup acknowledges it.
func (n *Node) mirror(ctx context.Context, v string) error {
	conn, err := n.env.Dial(ctx, n.backup+":"+strconv.Itoa(Port))
	if err != nil {
		return fmt.Errorf("mirrored: dial the backup: %w", err)
	}
	defer func() { _ = conn.Close() }()
	if _, err = io.WriteString(conn, v+"\n"); err != nil {
		return fmt.Errorf("mirrored: send the value: %w", err)
	}
	line, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return fmt.Errorf("mirrored: read the acknowledgement: %w", err)
	}
	if line != ack {
		return fmt.Errorf("mirrored: the backup replied %q", line)
	}
	return nil
}

// resend sends the value of the primary to the backup every Period, until
// the lifetime of the incarnation ends.
func (n *Node) resend() {
	life := n.env.Lifetime()
	for {
		n.env.Clock().Sleep(Period)
		if life.Err() != nil {
			return
		}
		if v, err := n.Get(life); err == nil {
			_ = n.mirror(life, v)
		}
	}
}

// accept serves each stream of the primary on the listener l, until the
// crash of the incarnation closes l.
func (n *Node) accept(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		n.env.Go(func() { n.serve(conn) })
	}
}

// serve stores each value of the stream conn and acknowledges it.
func (n *Node) serve(conn net.Conn) {
	defer func() { _ = conn.Close() }()
	r := bufio.NewReader(conn)
	for {
		line, err := r.ReadString('\n')
		if err != nil || n.store(strings.TrimSuffix(line, "\n")) != nil {
			return
		}
		if _, err := io.WriteString(conn, ack); err != nil {
			return
		}
	}
}

// store makes v the value of the node once v is durable. Each store writes a
// file of its own and renames it to the file of the value, so concurrent
// stores need no lock.
func (n *Node) store(v string) error {
	disk := n.env.Disk()
	name := valueFile + "." + strconv.FormatInt(n.writes.Add(1), 10)
	f, err := disk.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("mirrored: create the value: %w", err)
	}
	_, err = io.WriteString(f, v)
	if err == nil {
		err = f.Sync()
	}
	if err = errors.Join(err, f.Close()); err != nil {
		return fmt.Errorf("mirrored: write the value: %w", err)
	}
	if err := disk.Rename(name, valueFile); err != nil {
		return fmt.Errorf("mirrored: rename the value: %w", err)
	}
	if err := disk.SyncDir("."); err != nil {
		return fmt.Errorf("mirrored: sync the value: %w", err)
	}
	return nil
}
