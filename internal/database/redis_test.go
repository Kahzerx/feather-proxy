package database

import (
	"context"
	"testing"
)

func TestRedisConnection(t *testing.T) {
	client := NewRedisClient()
	if client.Ping(context.Background()).Err() != nil {
		t.Fatal("Unable to connect to redis")
	}
}
