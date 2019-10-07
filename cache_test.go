package cache

import (
	"testing"
	"time"
)

var testKey = "testKey"

func TestLocalCache_Get(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		fake *LocalCache
		want bool
	}{
		{
			name: "exist cache",
			fake: &LocalCache{Data: map[string][]byte{testKey: []byte("hoge")}, Expires: now.Add(60 * time.Second).Unix()},
			want: true,
		},
		{
			name: "cache expired",
			fake: &LocalCache{Data: map[string][]byte{testKey: []byte("hoge")}, Expires: now.Add(-60 * time.Second).Unix()},
			want: false,
		},
		{
			name: "cache not exist",
			fake: &LocalCache{Data: nil, Expires: now.Add(-60 * time.Second).Unix()},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fake.Get(testKey)
			if (got != nil) != tt.want {
				t.Fatalf("failed to get cache. want is %v, but got != nil is %v", tt.want, got != nil)
			}
		})
	}
}

func TestLocalCache_Set(t *testing.T) {
	tests := []struct {
		name string
		arg  []byte
		want bool
	}{
		{name: "success to set cache", arg: []byte("hoge"), want: false},
		{name: "failed to set for existing no cache", arg: nil, want: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := NewLocalCache(time.Now().Add(60 * time.Second).Unix())
			err := c.Set(testKey, tt.arg)
			if (err != nil) != tt.want {
				t.Fatalf("failed to set cache. err is %v but wantErr is %v", err, tt.want)
			}
		})
	}
}
