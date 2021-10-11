package notify

import (
	"errors"
	log "github.com/sirupsen/logrus"
)

type notifyFactoryImpl struct {
	options *Options
	hook    log.Hook
}

func (n *notifyFactoryImpl) Face() string {
	if n == nil || n.options == nil {
		return ""
	}
	return n.options.Name
}

func (n *notifyFactoryImpl) Create(args ...interface{}) (log.Hook, error) {
	var argc = len(args)
	if argc == 0 {
		if n.hook != nil {
			return n.hook, nil
		}
		args = append(args, n.options)
	}
	if n.options == args[0] {
		n.hook = NewHttpWebHook(*n.options)
		return n.hook, nil
	}
	var (
		info    = args[0]
		options = NewOptions(info)
	)
	if options == nil {
		return nil, errors.New("options missing call notifyFactoryImpl.Create")
	}
	return NewHttpWebHook(*options), nil
}

func CreateNotifyFactory(options *Options) *notifyFactoryImpl {
	if options == nil {
		return nil
	}
	var factory = new(notifyFactoryImpl)
	factory.options = options
	return factory
}
