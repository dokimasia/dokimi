// Package listener declares an interface whose double would have two
// members of one name.
package listener

// Listener has the method OnClose, which is the name of the field of the
// callable of Close.
//
//dokimi:stub
type Listener interface {
	// Close closes the listener.
	Close()
	// OnClose handles a close.
	OnClose()
}
