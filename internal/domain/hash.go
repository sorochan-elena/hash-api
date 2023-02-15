package domain

import (
	"errors"
	"time"
)

var (
	ErrNilHash = errors.New("nil hash")
	ErrExpired = errors.New("hash expired")
)

type Hash struct {
	Content   string
	CreatedAt time.Time
	expiresAt time.Time
}

func NewHash(content string, createdAt, expiresAt time.Time) Hash {
	return Hash{
		Content:   content,
		CreatedAt: createdAt,
		expiresAt: expiresAt,
	}
}

func (h Hash) Expired() bool {
	return time.Now().After(h.expiresAt)
}

func (h Hash) ExpiresAt() time.Time {
	return h.expiresAt
}

func (h Hash) Valid() bool {
	return h.Content != "" && !h.Expired()
}

type HashFunc = func(tm time.Time, ttl time.Duration) Hash
