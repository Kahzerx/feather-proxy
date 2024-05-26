package database

import (
	"context"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"testing"
	"time"
)

func TestRedisConnection(t *testing.T) {
	client := NewRedisClient()
	if client.Ping(context.Background()).Err() != nil {
		t.Fatal("Unable to connect to redis")
	}
	rs := redsync.New(goredis.NewPool(client))
	mutex := rs.NewMutex("feather-publishing", redsync.WithExpiry(10*time.Second))
	err := mutex.Lock()
	if err != nil {
		t.Fatal("Unable to create redsync Mutex, err: ", err.Error())
	}
	unlocked, err := mutex.Unlock()
	if err != nil {
		t.Fatal("Unable to unlock redsync Mutex, err: ", err.Error())
	}
	if !unlocked {
		t.Fatal("Unable to unlock redsync Mutex, lock status: ", unlocked)
	}
}
