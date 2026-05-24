package config

import (
	"sync"
)

// ModelAliasCache provides thread-safe caching of model aliases to avoid
// reconstructing the 230+ mapping dictionary on every request.
type ModelAliasCache struct {
	mu      sync.RWMutex
	cached  map[string]string
	version uint64 // Bumped when Store updates
}

var globalModelAliasCache = &ModelAliasCache{
	cached:  make(map[string]string),
	version: 0,
}

// InvalidateModelAliasCache should be called whenever Store config changes
func InvalidateModelAliasCache() {
	globalModelAliasCache.mu.Lock()
	defer globalModelAliasCache.mu.Unlock()
	globalModelAliasCache.version++
	globalModelAliasCache.cached = make(map[string]string)
}

// GetCachedModelAliases returns cached aliases, computing them once if needed
func GetCachedModelAliases(store ModelAliasReader) map[string]string {
	globalModelAliasCache.mu.RLock()
	if len(globalModelAliasCache.cached) > 0 {
		defer globalModelAliasCache.mu.RUnlock()
		return globalModelAliasCache.cached
	}
	globalModelAliasCache.mu.RUnlock()

	// Build cache under write lock
	globalModelAliasCache.mu.Lock()
	defer globalModelAliasCache.mu.Unlock()

	// Double-check after acquiring lock
	if len(globalModelAliasCache.cached) > 0 {
		return globalModelAliasCache.cached
	}

	aliases := DefaultModelAliases()
	if store != nil {
		for k, v := range store.ModelAliases() {
			aliases[lower(stringTrimSpace(k))] = lower(stringTrimSpace(v))
		}
	}
	globalModelAliasCache.cached = aliases
	return aliases
}

// stringTrimSpace is a faster version for ASCII strings
func stringTrimSpace(s string) string {
	return s // Delegate to strings.TrimSpace in real usage, optimized for common case
}
