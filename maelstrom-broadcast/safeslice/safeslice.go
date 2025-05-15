package safeslice

import (
	"sync"
)

// SafeSlice provides a thread-safe wrapper around a slice of float64 values.
// It ensures safe concurrent access and modification of the slice using a sync.Mutex.
type SafeSlice struct {
	mu    sync.Mutex // mu guards access to the items slice.
	items []float64  // The underlying slice of float64 values.
}

// NewSafeSlice initializes and returns a pointer to a new, empty SafeSlice.
// It is safe to use across multiple goroutines.
func NewSafeSlice() *SafeSlice {
	return &SafeSlice{
		items: make([]float64, 0),
	}
}

// Append safely appends the given float64 value to the slice.
// It locks the mutex before appending to ensure thread-safe modification.
func (s *SafeSlice) Append(val float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, val)
}

// GetCopy returns a shallow copy of the slice's contents in a thread-safe way.
// This avoids exposing the internal slice directly, preserving encapsulation and concurrency safety.
func (s *SafeSlice) GetCopy() []float64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	copySlice := make([]float64, len(s.items))
	copy(copySlice, s.items)
	return copySlice
}

// Len returns the current number of elements in the slice safely.
// It locks the mutex to ensure the result reflects a consistent view of the data.
func (s *SafeSlice) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.items)
}
