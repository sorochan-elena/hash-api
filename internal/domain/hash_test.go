package domain_test

import (
	"github.com/stretchr/testify/assert"
	"hash-api/internal/domain"
	"testing"
	"time"
)

func TestHashValid(t *testing.T) {
	tests := []struct {
		name string
		hash domain.Hash
		want bool
	}{
		{
			name: "empty",
			hash: domain.Hash{},
			want: false,
		},
		{
			name: "empty content",
			hash: domain.Hash{CreatedAt: time.Now()},
			want: false,
		},
		{
			name: "empty time",
			hash: domain.Hash{Content: "aaa"},
			want: false,
		},
		{
			name: "expired",
			hash: domain.NewHash("aaa", time.Now().Add(time.Second), time.Now()),
			want: false,
		},
		{
			name: "valid",
			hash: domain.NewHash("aaa", time.Now(), time.Now().Add(time.Second)),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.hash.Valid())
		})
	}
}
