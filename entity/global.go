package entity

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type (

	HookFactory func(args ...[]byte) (logrus.Hook, error)

	hookProvider struct {
		locker  sync.RWMutex
		drivers map[string]HookFactory
	}
)

func CreateProvider() *hookProvider {
	var provider = new(hookProvider)
	provider.locker = sync.RWMutex{}
	provider.drivers = make(map[string]HookFactory)
	return provider
}

func (hook *hookProvider) Register(key string, factory HookFactory) bool {
	hook.locker.Lock()
	defer hook.locker.Unlock()
	if v, ok := hook.drivers[key]; ok && v != nil {
		return false
	}
	hook.drivers[key] = factory
	return true
}

func (hook *hookProvider) Remove(key string) bool {
	hook.locker.Lock()
	defer hook.locker.Unlock()
	if _, ok := hook.drivers[key]; ok {
		delete(hook.drivers, key)
	}
	return true
}

func (hook *hookProvider) Len() int {
	hook.locker.Lock()
	defer hook.locker.Unlock()
	return len(hook.drivers)
}

func (hook *hookProvider) Exists(key string) bool {
	hook.locker.Lock()
	defer hook.locker.Unlock()
	if _, ok := hook.drivers[key]; ok {
		return true
	}
	return false
}

func (hook *hookProvider) Get(key string) (HookFactory, bool) {
	hook.locker.Lock()
	defer hook.locker.Unlock()
	var factory, ok = hook.drivers[key]
	return factory, ok
}

func (hook *hookProvider) Resolve(key string, args ...[]byte) (logrus.Hook, error) {
	var factory, ok = hook.Get(key)
	if !ok {
		return nil, ErrNotExists
	}
	return factory(args...)
}
