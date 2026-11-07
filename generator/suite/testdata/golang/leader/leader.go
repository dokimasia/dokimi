// Package leader declares the client of a leader, whose suite the directive
// out routes into the kit leadertest.
package leader

// Primary is the client of the replication calls of a leader. The test of its
// suite takes the options of its falsification from the external test package
// of the kit.
//
//dokimi:suite falsify=falsifyOptions
//dokimi:out path=leadertest/
type Primary interface {
	// Commit returns the index of the last committed entry.
	//
	//dokimi:budget allocs=0
	Commit() int
}
