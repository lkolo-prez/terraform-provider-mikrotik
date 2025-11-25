package client

import (
	"fmt"
	"sync"
)

// BgpBatchOperations provides optimized batch operations for BGP resources
// to reduce RouterOS API calls and improve performance
type BgpBatchOperations struct {
	client *Mikrotik
	cache  *BgpCache
}

// BgpCache provides in-memory caching for BGP resources to reduce API calls
type BgpCache struct {
	mu                sync.RWMutex
	instancesV7       map[string]*BgpInstanceV7
	connections       map[string]*BgpConnection
	templates         map[string]*BgpTemplate
	sessions          map[string]*BgpSession
	instancesV7Valid  bool
	connectionsValid  bool
	templatesValid    bool
	sessionsValid     bool
}

// NewBgpBatchOperations creates a new batch operations handler
func (c *Mikrotik) NewBgpBatchOperations() *BgpBatchOperations {
	return &BgpBatchOperations{
		client: c,
		cache: &BgpCache{
			instancesV7: make(map[string]*BgpInstanceV7),
			connections: make(map[string]*BgpConnection),
			templates:   make(map[string]*BgpTemplate),
			sessions:    make(map[string]*BgpSession),
		},
	}
}

// InvalidateCache clears all cached data
func (b *BgpBatchOperations) InvalidateCache() {
	b.cache.mu.Lock()
	defer b.cache.mu.Unlock()
	
	b.cache.instancesV7Valid = false
	b.cache.connectionsValid = false
	b.cache.templatesValid = false
	b.cache.sessionsValid = false
}

// GetOrFetchInstanceV7 retrieves BGP instance from cache or fetches from RouterOS
func (b *BgpBatchOperations) GetOrFetchInstanceV7(name string) (*BgpInstanceV7, error) {
	b.cache.mu.RLock()
	if b.cache.instancesV7Valid {
		if instance, ok := b.cache.instancesV7[name]; ok {
			b.cache.mu.RUnlock()
			return instance, nil
		}
	}
	b.cache.mu.RUnlock()

	// Fetch from RouterOS
	instance, err := b.client.FindBgpInstanceV7(name)
	if err != nil {
		return nil, err
	}

	// Update cache
	b.cache.mu.Lock()
	b.cache.instancesV7[name] = instance
	b.cache.instancesV7Valid = true
	b.cache.mu.Unlock()

	return instance, nil
}

// GetOrFetchConnection retrieves BGP connection from cache or fetches from RouterOS
func (b *BgpBatchOperations) GetOrFetchConnection(name string) (*BgpConnection, error) {
	b.cache.mu.RLock()
	if b.cache.connectionsValid {
		if conn, ok := b.cache.connections[name]; ok {
			b.cache.mu.RUnlock()
			return conn, nil
		}
	}
	b.cache.mu.RUnlock()

	// Fetch from RouterOS
	conn, err := b.client.FindBgpConnection(name)
	if err != nil {
		return nil, err
	}

	// Update cache
	b.cache.mu.Lock()
	b.cache.connections[name] = conn
	b.cache.connectionsValid = true
	b.cache.mu.Unlock()

	return conn, nil
}

// GetOrFetchTemplate retrieves BGP template from cache or fetches from RouterOS
func (b *BgpBatchOperations) GetOrFetchTemplate(name string) (*BgpTemplate, error) {
	b.cache.mu.RLock()
	if b.cache.templatesValid {
		if tpl, ok := b.cache.templates[name]; ok {
			b.cache.mu.RUnlock()
			return tpl, nil
		}
	}
	b.cache.mu.RUnlock()

	// Fetch from RouterOS
	tpl, err := b.client.FindBgpTemplate(name)
	if err != nil {
		return nil, err
	}

	// Update cache
	b.cache.mu.Lock()
	b.cache.templates[name] = tpl
	b.cache.templatesValid = true
	b.cache.mu.Unlock()

	return tpl, nil
}

// GetOrFetchSession retrieves BGP session from cache or fetches from RouterOS
func (b *BgpBatchOperations) GetOrFetchSession(name string) (*BgpSession, error) {
	b.cache.mu.RLock()
	if b.cache.sessionsValid {
		if session, ok := b.cache.sessions[name]; ok {
			b.cache.mu.RUnlock()
			return session, nil
		}
	}
	b.cache.mu.RUnlock()

	// Fetch from RouterOS
	session, err := b.client.FindBgpSession(name)
	if err != nil {
		return nil, err
	}

	// Update cache
	b.cache.mu.Lock()
	b.cache.sessions[name] = session
	b.cache.sessionsValid = true
	b.cache.mu.Unlock()

	return session, nil
}

// BatchAddConnections adds multiple BGP connections in optimized manner
func (b *BgpBatchOperations) BatchAddConnections(connections []*BgpConnection) ([]*BgpConnection, []error) {
	results := make([]*BgpConnection, len(connections))
	errors := make([]error, len(connections))

	for i, conn := range connections {
		result, err := b.client.AddBgpConnection(conn)
		results[i] = result
		errors[i] = err
		
		// Update cache on success
		if err == nil && result != nil {
			b.cache.mu.Lock()
			b.cache.connections[result.Name] = result
			b.cache.mu.Unlock()
		}
	}

	return results, errors
}

// BatchUpdateConnections updates multiple BGP connections in optimized manner
func (b *BgpBatchOperations) BatchUpdateConnections(connections []*BgpConnection) ([]*BgpConnection, []error) {
	results := make([]*BgpConnection, len(connections))
	errors := make([]error, len(connections))

	for i, conn := range connections {
		result, err := b.client.UpdateBgpConnection(conn)
		results[i] = result
		errors[i] = err
		
		// Update cache on success
		if err == nil && result != nil {
			b.cache.mu.Lock()
			b.cache.connections[result.Name] = result
			b.cache.mu.Unlock()
		}
	}

	return results, errors
}

// PreloadAllSessions fetches all BGP sessions at once for caching
// This is useful when you need to query multiple sessions
func (b *BgpBatchOperations) PreloadAllSessions() error {
	sessions, err := b.client.ListBgpSessions()
	if err != nil {
		return fmt.Errorf("failed to preload sessions: %w", err)
	}

	b.cache.mu.Lock()
	defer b.cache.mu.Unlock()

	b.cache.sessions = make(map[string]*BgpSession, len(sessions))
	for _, session := range sessions {
		if session != nil {
			b.cache.sessions[session.Name] = session
		}
	}
	b.cache.sessionsValid = true

	return nil
}

// GetCacheStats returns cache statistics for monitoring
func (b *BgpBatchOperations) GetCacheStats() map[string]interface{} {
	b.cache.mu.RLock()
	defer b.cache.mu.RUnlock()

	return map[string]interface{}{
		"instances_v7_count": len(b.cache.instancesV7),
		"connections_count":  len(b.cache.connections),
		"templates_count":    len(b.cache.templates),
		"sessions_count":     len(b.cache.sessions),
		"instances_valid":    b.cache.instancesV7Valid,
		"connections_valid":  b.cache.connectionsValid,
		"templates_valid":    b.cache.templatesValid,
		"sessions_valid":     b.cache.sessionsValid,
	}
}
