package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	dirPermission  = 0755
	filePermission = 0600
)

var (
	ErrCacheMiss    = errors.New("cache miss")
	ErrCacheExpired = errors.New("cache entry expired")
	ErrCacheCorrupt = errors.New("cache entry corrupted")
)

type FileCache struct {
	baseDir string
	ttl     time.Duration
}

func New(baseDir string, ttl time.Duration) *FileCache {
	return &FileCache{
		baseDir: baseDir,
		ttl:     ttl,
	}
}

type Entry struct {
	Data      json.RawMessage `json:"data"`
	ExpiresAt time.Time       `json:"expires_at"`
}

func (c *FileCache) key(ecosystem, name, version string) string {
	safeName := strings.ReplaceAll(name, "/", "_")
	safeName = strings.ReplaceAll(safeName, ":", "_")
	return filepath.Join(c.baseDir, ecosystem, fmt.Sprintf("%s@%s.json", safeName, version))
}

func (c *FileCache) Get(ecosystem, name, version string) ([]byte, error) {
	path := c.key(ecosystem, name, version)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrCacheMiss
		}
		return nil, fmt.Errorf("read cache: %w", err)
	}

	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		if err := os.Remove(path); err != nil {
			return nil, fmt.Errorf("failed to remove corrupted cache: %w", err)
		}
		return nil, ErrCacheCorrupt
	}

	if time.Now().After(entry.ExpiresAt) {
		if err := os.Remove(path); err != nil {
			return nil, fmt.Errorf("failed to remove expired cache: %w", err)
		}
		return nil, ErrCacheExpired
	}

	return entry.Data, nil
}

func (c *FileCache) Set(ecosystem, name, version string, data []byte) error {
	path := c.key(ecosystem, name, version)
	if err := os.MkdirAll(filepath.Dir(path), dirPermission); err != nil {
		return fmt.Errorf("create cache dir: %w", err)
	}

	entry := Entry{
		Data:      data,
		ExpiresAt: time.Now().Add(c.ttl),
	}
	raw, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshal cache entry: %w", err)
	}

	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, raw, filePermission); err != nil {
		return fmt.Errorf("write cache: %w", err)
	}
	return os.Rename(tmpPath, path)
}
