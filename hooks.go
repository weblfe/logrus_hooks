package logrus_hooks

import (
	"github.com/sirupsen/logrus"
	"github.com/weblfe/logrus_hooks/entity"
)

var (
	hooksProvider = entity.CreateProvider()
)

func Register(hook string, factory entity.HookFactory) bool {
	return hooksProvider.Register(hook, factory)
}

func Resolve(hook string, args ...[]byte) (logrus.Hook, error) {
	return hooksProvider.Resolve(hook, args...)
}

func Exists(hook string) bool {
	return hooksProvider.Exists(hook)
}
