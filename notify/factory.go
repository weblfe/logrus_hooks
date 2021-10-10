package notify

import (
		log "github.com/sirupsen/logrus"
		"github.com/weblfe/logrus_hooks/faced"
)

type notifyMgrImpl struct {
		factories map[string]faced.HookFactory
}

var (
		notifyMgr = newNotifyMgrImpl()
)

func newNotifyMgrImpl() *notifyMgrImpl  {
		var mgr = new(notifyMgrImpl)
		mgr.factories = make(map[string]faced.HookFactory)
		return mgr
}

func (mgr *notifyMgrImpl)invoke(name string) faced.Creator  {
		return func(args...interface{}) (log.Hook,error) {
				return nil,nil
		}
}

func (mgr *notifyMgrImpl)Add(factory faced.HookFactory) *notifyMgrImpl {
		if factory == nil {
				return mgr
		}
		if _,ok:=mgr.factories[factory.Face()];ok {
				return mgr
		}
		mgr.factories[factory.Face()] = factory
		return mgr
}

func (mgr *notifyMgrImpl)Hooks() []string  {
		return nil
}

func (mgr *notifyMgrImpl)Register(hookMgr faced.HookMgr)  {
			if hookMgr==nil {
					return
			}
		for _,v:=range mgr.Hooks() {
				hookMgr.Register(v,mgr.invoke(v))
		}
}

func Register(mgr faced.HookMgr)  {
		notifyMgr.Register(mgr)
}