package cache

import (
	"sync"
	"time"

	"github.com/sixworks/go-snmp-prometheus-getter/internal/elasticsearch"
)

// ConfigCache provides a thread-safe cache for device configurations
type ConfigCache struct {
	configs  map[string]elasticsearch.Config
	mu       sync.RWMutex
	ttl      time.Duration
	updated  time.Time
}

// NewConfigCache creates a new configuration cache with the specified TTL
func NewConfigCache(ttl time.Duration) *ConfigCache {
	return &ConfigCache{
		configs: make(map[string]elasticsearch.Config),
		ttl:     ttl,
	}
}

// Get retrieves a configuration by ID
func (c *ConfigCache) Get(id string) (elasticsearch.Config, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	config, exists := c.configs[id]
	return config, exists
}

// GetAll returns all cached configurations
func (c *ConfigCache) GetAll() []elasticsearch.Config {
	c.mu.RLock()
	defer c.mu.RUnlock()

	configs := make([]elasticsearch.Config, 0, len(c.configs))
	for _, config := range c.configs {
		configs = append(configs, config)
	}
	return configs
}

// Set stores a configuration in the cache
func (c *ConfigCache) Set(config elasticsearch.Config) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.configs[config.ID] = config
}

// SetAll replaces all configurations in the cache
func (c *ConfigCache) SetAll(configs []elasticsearch.Config) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.configs = make(map[string]elasticsearch.Config, len(configs))
	for _, config := range configs {
		c.configs[config.ID] = config
	}
	c.updated = time.Now()
}

// Delete removes a configuration from the cache
func (c *ConfigCache) Delete(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.configs, id)
}

// Clear removes all configurations from the cache
func (c *ConfigCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.configs = make(map[string]elasticsearch.Config)
	c.updated = time.Time{}
}

// IsExpired checks if the cache has expired
func (c *ConfigCache) IsExpired() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return time.Since(c.updated) > c.ttl
}

// LastUpdated returns when the cache was last updated
func (c *ConfigCache) LastUpdated() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.updated
}

// Count returns the number of configurations in the cache
func (c *ConfigCache) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.configs)
}
