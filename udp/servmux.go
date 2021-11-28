package udplevis

import (
	"levis"
	"net"
)

// ServeMux provides mappings from a common endpoint to handlers by request path.
type ServeMux struct {
	m map[string]muxEntry
}

type muxEntry struct {
	h       Handler
	pattern string
}

// NewServeMux creates a new ServeMux.
func NewServeMux() *ServeMux { return &ServeMux{m: make(map[string]muxEntry)} }

// Does path match pattern?
func pathMatch(pattern, path string) bool {
	if len(pattern) == 0 {
		// should not happen
		return false
	}
	n := len(pattern)
	if pattern[n-1] != '/' {
		return pattern == path
	}
	return len(path) >= n && path[0:n] == pattern
}

// Find a handler on a handler map given a path string
// Most-specific (longest) pattern wins
func (mux *ServeMux) match(path string) (h Handler, pattern string) {
	var n = 0
	for k, v := range mux.m {
		if !pathMatch(k, path) {
			continue
		}
		if h == nil || len(k) > n {
			n = len(k)
			h = v.h
			pattern = v.pattern
		}
	}
	return
}

var _ = Handler(&ServeMux{})

// ServeLevis handles a single Levis message.  The message arrives from
// the given listener having originated from the given UDPAddr.
func (mux *ServeMux) ServeLevis(l *net.UDPConn, a *net.UDPAddr, m *levis.Message) *levis.Message {
	h, _ := mux.match(levis.ProtocolName)
	return h.ServeLevis(l, a, m)
}

// Handle configures a handler for the given path.
func (mux *ServeMux) Handle(handler Handler) {
	pattern := levis.ProtocolName
	mux.m[pattern] = muxEntry{h: handler, pattern: pattern}
}

// HandleFunc configures a handler for the given path.
func (mux *ServeMux) HandleFunc(f func(l *net.UDPConn, a *net.UDPAddr, m *levis.Message) *levis.Message) {
	mux.Handle(FuncHandler(f))
}
