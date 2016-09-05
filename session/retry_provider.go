package session

import (
	"time"

	"github.com/cenkalti/backoff"
	ms "github.com/go-macaron/session"
)

type retryProvider struct {
	provider ms.Provider
}

func NewRetryProvider(provider ms.Provider) ms.Provider {
	return &retryProvider{
		provider: provider,
	}
}

func (p *retryProvider) backoff() backoff.BackOff {
	result := backoff.NewExponentialBackOff()
	result.InitialInterval = time.Millisecond * 100
	result.MaxElapsedTime = time.Second * 10
	result.Reset()
	return result
}

// Init initializes session provider.
func (p *retryProvider) Init(gclifetime int64, config string) error {
	op := func() error {
		return p.provider.Init(gclifetime, config)
	}
	return maskAny(backoff.Retry(op, p.backoff()))
}

// Read returns raw session store by session ID.
func (p *retryProvider) Read(sid string) (ms.RawStore, error) {
	var result ms.RawStore
	op := func() error {
		var err error
		result, err = p.provider.Read(sid)
		return err
	}
	err := maskAny(backoff.Retry(op, p.backoff()))
	if err != nil {
		return result, maskAny(err)
	}
	return NewRetryStore(result), nil
}

// Exist returns true if session with given ID exists.
func (p *retryProvider) Exist(sid string) bool {
	return p.provider.Exist(sid)
}

// Destory deletes a session by session ID.
func (p *retryProvider) Destory(sid string) error {
	op := func() error {
		return p.provider.Destory(sid)
	}
	return maskAny(backoff.Retry(op, p.backoff()))

}

// Regenerate regenerates a session store from old session ID to new one.
func (p *retryProvider) Regenerate(oldsid, sid string) (ms.RawStore, error) {
	var result ms.RawStore
	op := func() error {
		var err error
		result, err = p.provider.Regenerate(oldsid, sid)
		return err
	}
	err := maskAny(backoff.Retry(op, p.backoff()))
	if err != nil {
		return result, maskAny(err)
	}
	return NewRetryStore(result), nil
}

// Count counts and returns number of sessions.
func (p *retryProvider) Count() int {
	return p.provider.Count()
}

// GC calls GC to clean expired sessions.
func (p *retryProvider) GC() {
	p.provider.GC()
}
