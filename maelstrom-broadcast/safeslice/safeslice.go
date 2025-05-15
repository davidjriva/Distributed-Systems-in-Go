package safeslice

import (
	"sync"
)

// SafeSlice is a thread-safe slice of float64 values.
type SafeSlice struct {
	mu    sync.Mutex
	items []float64
}

// NewSafeSlice creates and returns a new SafeSlice.
func NewSafeSlice() *SafeSlice {
	return &SafeSlice{
		items: make([]float64, 0),
	}
}

// Append adds a new value to the slice in a thread-safe manner.
func (s *SafeSlice) Append(val float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, val)
}

// GetCopy returns a thread-safe copy of the underlying slice.
func (s *SafeSlice) GetCopy() []float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	copySlice := make([]float64, len(s.items))
	copy(copySlice, s.items)
	return copySlice
}

// Len returns the length of the slice in a thread-safe manner.
func (s *SafeSlice) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.items)
}
