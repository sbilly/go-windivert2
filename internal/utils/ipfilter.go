package utils

import (
	"net"
	"sync"
)

// IPFilter represents an IP filter
type IPFilter struct {
	sync.RWMutex
	tree *IPTree
}

// NewIPFilter creates a new IP filter
func NewIPFilter() *IPFilter {
	return &IPFilter{
		tree: NewIPTree(),
	}
}

// Add adds an IP address to the filter
func (f *IPFilter) Add(ip net.IP) {
	f.Lock()
	defer f.Unlock()
	f.tree.Insert(ip)
}

// Lookup checks if an IP address exists in the filter
func (f *IPFilter) Lookup(ip net.IP) bool {
	f.RLock()
	defer f.RUnlock()
	return f.tree.Lookup(ip)
}
