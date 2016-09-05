package session

import (
	"time"

	"github.com/cenkalti/backoff"
	ms "github.com/go-macaron/session"
)

type retryStore struct {
	store ms.RawStore
}

func NewRetryStore(store ms.RawStore) ms.RawStore {
	return &retryStore{
		store: store,
	}
}

func (rs *retryStore) backoff() backoff.BackOff {
	result := backoff.NewExponentialBackOff()
	result.InitialInterval = time.Millisecond * 100
	result.MaxElapsedTime = time.Second * 10
	result.Reset()
	return result
}

// Set sets value to given key in session.
func (rs *retryStore) Set(x interface{}, y interface{}) error {
	op := func() error {
		return rs.store.Set(x, y)
	}
	return maskAny(backoff.Retry(op, rs.backoff()))
}

// Get gets value by given key in session.
func (rs *retryStore) Get(key interface{}) interface{} {
	return rs.store.Get(key)
}

// Delete deletes a key from session.
func (rs *retryStore) Delete(key interface{}) error {
	op := func() error {
		return rs.store.Delete(key)
	}
	return maskAny(backoff.Retry(op, rs.backoff()))
}

// ID returns current session ID.
func (rs *retryStore) ID() string {
	return rs.store.ID()
}

// Release releases session resource and save data to provider.
func (rs *retryStore) Release() error {
	op := func() error {
		return rs.store.Release()
	}
	return maskAny(backoff.Retry(op, rs.backoff()))
}

// Flush deletes all session data.
func (rs *retryStore) Flush() error {
	op := func() error {
		return rs.store.Flush()
	}
	return maskAny(backoff.Retry(op, rs.backoff()))
}
