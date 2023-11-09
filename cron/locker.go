package cron

import (
	"context"
	"errors"
)

var (
	ErrFailedToConnectToRedis = errors.New("cron: failed to connect to redis")
	ErrFailedToObtainLock     = errors.New("cron: failed to obtain lock")
	ErrFailedToReleaseLock    = errors.New("cron: failed to release lock")
)

// Locker represents the required interface to lock jobs when running multiple schedulers.
type Locker interface {
	// Lock if an error is returned by lock, the job will not be scheduled.
	Lock(ctx context.Context, key string) (Lock, error)
}

// Lock represents an obtained lock
type Lock interface {
	Unlock(ctx context.Context) error
}

// Elector determines the leader from instances asking to be the leader. Only
// the leader runs jobs. If the leader goes down, a new leader will be elected.
type Elector interface {
	// IsLeader should return an error if the job should not be scheduled and nil if the job should be scheduled.
	IsLeader(ctx context.Context) error
}
