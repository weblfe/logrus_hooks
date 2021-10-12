package notify

import (
	"errors"
	"github.com/weblfe/logrus_hooks/facede"
)

type notifyMgrImpl struct {
	factories map[string]facede.HookFactory
}

var (
	notifyMgr = newNotifyMgrImpl()
)

func newNotifyMgrImpl() *notifyMgrImpl {
	var mgr = new(notifyMgrImpl)
	mgr.factories = make(map[string]facede.HookFactory)
	return mgr
}

func (mgr *notifyMgrImpl) invoke(name string) facede.Creator {
	if name == "" || mgr == nil {
		return nil
	}
	// 1. 注册中查找 构造
	if v, ok := mgr.factories[name]; ok {
		return v.Create
	}
	// 2. 环境变中注册参数构造
	var creator = mgr.findByEnv(name)
	if creator != nil {
		return creator
	}
	// 3. nil + empty name
	return nil
}

func (mgr *notifyMgrImpl) findByEnv(faceID string) facede.Creator {
	if faceID == "" {
		return nil
	}
	var options = NewOptionWithEnvPrefix(faceID)
	if options == nil || options.Url == "" {
		return nil
	}
	var factory = CreateNotifyFactory(options)
	// 注册缓存
	if err := mgr.Add(factory); err != nil {
		return nil
	}
	return factory.Create
}

func (mgr *notifyMgrImpl) Add(factory facede.HookFactory) error {
	if factory == nil {
		return errors.New("required not nil factory")
	}
	var face = factory.Face()
	// miss face
	if face == "" {
		return errors.New("factory miss face id")
	}
	if _, ok := mgr.factories[face]; ok {
		return nil
	}
	mgr.factories[face] = factory
	return nil
}

func (mgr *notifyMgrImpl) Hooks() []string {
	if mgr == nil || len(mgr.factories) <= 0 {
		return nil
	}
	var hooks []string
	for k := range mgr.factories {
		hooks = append(hooks, k)
	}
	return hooks
}

func (mgr *notifyMgrImpl) Register(hookMgr facede.HookMgr) {
	if hookMgr == nil {
		return
	}
	for _, v := range mgr.Hooks() {
		hookMgr.Register(v, mgr.invoke(v))
	}
}

func Register(mgr facede.HookMgr) {
	notifyMgr.Register(mgr)
}
