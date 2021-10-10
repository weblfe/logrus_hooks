package entity

import (
	"github.com/sirupsen/logrus"
		"github.com/weblfe/logrus_hooks/faced"
		"sync"
)

type (
		// hookMgrImpl 注册提供者
	hookMgrImpl struct {
		locker  sync.RWMutex
		drivers map[string]faced.Creator
	}
)

func CreateProvider() *hookMgrImpl {
	var provider = new(hookMgrImpl)
	provider.locker = sync.RWMutex{}
	provider.drivers = make(map[string]faced.Creator)
	return provider
}

func (hook *hookMgrImpl) Register(key string, factory faced.Creator) bool {
	hook.locker.Lock()
	defer hook.locker.Unlock()
	if _, ok := hook.drivers[key]; ok  {
		return false
	}
	hook.drivers[key] = factory
	return true
}

func (hook *hookMgrImpl) Remove(key string) bool {
	hook.locker.Lock()
	defer hook.locker.Unlock()
	if _, ok := hook.drivers[key]; ok {
		delete(hook.drivers, key)
	}
	return true
}

func (hook *hookMgrImpl) Len() int {
	hook.locker.Lock()
	defer hook.locker.Unlock()
	return len(hook.drivers)
}

func (hook *hookMgrImpl) Exists(key string) bool {
	hook.locker.Lock()
	defer hook.locker.Unlock()
	if _, ok := hook.drivers[key]; ok {
		return true
	}
	return false
}

func (hook *hookMgrImpl) Get(key string) (faced.Creator, bool) {
	hook.locker.Lock()
	defer hook.locker.Unlock()
	var factory, ok = hook.drivers[key]
	return factory, ok
}

func (hook *hookMgrImpl) Resolve(key string, args ...interface{}) (logrus.Hook, error) {
	var factory, ok = hook.Get(key)
	if !ok {
		return nil, ErrNotExists
	}
	return factory(args...)
}

// Replace 替换
func (hook *hookMgrImpl) Replace(key string, factory faced.Creator) bool {
		hook.locker.Lock()
		defer hook.locker.Unlock()
		hook.drivers[key] = factory
		return true
}
