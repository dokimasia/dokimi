// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package spec

import (
	"strings"

	"go.dokimi.dev/eidos/lang/naming"
)

// Key returns the key form of a check identity, which joins one check across
// languages. Each of the declaration and the callable becomes its words in
// lower case, joined by hyphens, and the segment and the case stay as they
// are: Store.ReadAll/zero-on-error has the key store.read-all/zero-on-error.
// An identity without a segment, such as Store.PutAll, is the identity of a
// callable, and its key is store.put-all.
func Key(identity string) string {
	head, rest, segmented := strings.Cut(identity, "/")
	declaration, callable, method := strings.Cut(head, ".")
	key := words(declaration)
	if method {
		key += "." + words(callable)
	}
	if segmented {
		key += "/" + rest
	}
	return key
}

// words returns the words of name in lower case, joined by hyphens.
func words(name string) string {
	return strings.ToLower(strings.Join(naming.Words(name), "-"))
}
