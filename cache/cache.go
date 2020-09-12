package cache

import (
	"time"
)

type Cache interface {
	Delete(key string) error
	Get(key string, obj interface{}) error
	Set(key string, obj interface{}, timeout time.Duration) error
	Close() error
}
