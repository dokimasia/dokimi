// Package remote calls the services of the shop. This file imports no
// package of errors, and link.go imports errs.
package remote

// Remote calls one service.
type Remote interface {
	// Call sends the request, and its double arms the timeout of errs.
	//
	//dokimi:fault errors=[errs.ErrTimeout]
	Call(request string) (string, error)
}
