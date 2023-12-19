package freelru

import (
	"sync"
	"testing"
)

const CAP = 4

// TestRaceCondition tests that the synced LRU is safe to use concurrently.
// Test with 'go test . -race'.
func TestRaceCondition(t *testing.T) {
	lru, err := NewSynced[uint64, int](CAP, hashUint64)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	wg := sync.WaitGroup{}

	call := func(fn func()) {
		wg.Add(1)
		go func() {
			fn()
			wg.Done()
		}()
	}

	call(func() { lru.SetLifetime(0) })
	call(func() { lru.SetOnEvict(nil) })
	call(func() { _ = lru.Len() })
	call(func() { _ = lru.AddWithLifetime(1, 1, 0) })
	call(func() { _ = lru.Add(1, 1) })
	call(func() { _, _ = lru.Get(1) })
	call(func() { _, _ = lru.Peek(1) })
	call(func() { _ = lru.Contains(1) })
	call(func() { _ = lru.Remove(1) })
	call(func() { _ = lru.Keys() })
	call(func() { lru.Purge() })
	call(func() { lru.dump() })
	call(func() { lru.PrintStats() })

	wg.Wait()
}
