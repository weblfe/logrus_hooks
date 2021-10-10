package faced

import log "github.com/sirupsen/logrus"

// HookFactory 钩子工厂
type (

	HookFactory interface {
		Face() string
		Create(args ...interface{}) (log.Hook, error)
	}
	// Creator hook 构造器
	Creator func(args ...interface{}) (log.Hook, error)

	// HookMgr hook 管理器
	HookMgr interface {
		Remove(key string) bool
		Exists(key string) bool
		Get(key string) (Creator, bool)
		Replace(key string, factory Creator) bool
		Register(key string, factory Creator) bool
		Resolve(key string, args ...interface{}) (log.Hook, error)
	}

)
