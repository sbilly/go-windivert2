package utils

// AppFilter represents an application filter
type AppFilter struct {
	pids map[uint32]struct{}
}

// NewAppFilter creates a new application filter
func NewAppFilter() *AppFilter {
	return &AppFilter{
		pids: make(map[uint32]struct{}),
	}
}

// Add adds a process ID to the filter
func (f *AppFilter) Add(pid uint32) {
	f.pids[pid] = struct{}{}
}

// Remove removes a process ID from the filter
func (f *AppFilter) Remove(pid uint32) {
	delete(f.pids, pid)
}

// Lookup checks if a process ID exists in the filter
func (f *AppFilter) Lookup(pid uint32) bool {
	_, ok := f.pids[pid]
	return ok
}
