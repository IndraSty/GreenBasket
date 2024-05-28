package domain

import "time"

type CacheRepository interface {
	Get(key string) ([]byte, error)
	Del(key string) error
	Set(key string, entry []byte, expiration time.Duration) error
}
