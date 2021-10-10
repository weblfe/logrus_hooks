package logrus_hooks

import (
	"github.com/sirupsen/logrus"
	"github.com/weblfe/logrus_hooks/entity"
	"github.com/weblfe/logrus_hooks/faced"
	"github.com/weblfe/logrus_hooks/notify"
	"github.com/weblfe/logrus_hooks/rotate"
)

var (
	hooksProvider = entity.CreateProvider()
)

func Register(hook string, factory faced.Creator) bool {
	return hooksProvider.Register(hook, factory)
}

func Resolve(hook string, args ...interface{}) (logrus.Hook, error) {
	return hooksProvider.Resolve(hook, args...)
}

func Exists(hook string) bool {
	return hooksProvider.Exists(hook)
}

// Add 添加创建工厂
func Add(factory faced.HookFactory) bool {
	if factory == nil {
		return false
	}
	return Register(factory.Face(), factory.Create)
}

func GetMgr() faced.HookMgr {
	return hooksProvider
}

func init() {
	// 被动注册
	Add(rotate.CreateRotateFactory())
	// 主动注册
	notify.Register(GetMgr())
}
