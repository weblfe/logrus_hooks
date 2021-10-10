package rotate

import (
		"github.com/lestrrat-go/file-rotatelogs"
		"github.com/rifflock/lfshook"
		log "github.com/sirupsen/logrus"
		"github.com/weblfe/logrus_hooks/entity"
		"github.com/weblfe/logrus_hooks/utils"
)

type (
	// 日志分割hook 工厂
	rotateHookFactory struct {
		defaultOption *Options
		name          string
	}
)

const (
	HookName = "rotate"
	LogExt   = ".log"
)

func CreateRotateFactory() *rotateHookFactory {
	var factory = new(rotateHookFactory)
	factory.name = HookName
	return factory
}


// Create 构建按日分割日志 hook
func (factory *rotateHookFactory) Create(args ...interface{}) (log.Hook, error) {
	if len(args) == 0 {
		return factory.newLfsHook(factory.getDefaultOption()), nil
	}
	var options = NewOption(args[0])
	return factory.newLfsHook(options), nil
}

func (factory *rotateHookFactory) getDefaultOption() *Options {
	if factory.defaultOption == nil {
		factory.defaultOption = CreateOptionsWithEnv(utils.UpperCase)
	}
	return factory.defaultOption
}

func (factory *rotateHookFactory) Face() string {
	return factory.name
}

func (factory *rotateHookFactory) SetDefaultOption(options *Options) *rotateHookFactory {
	if factory.defaultOption == nil && options != nil {
		factory.defaultOption = options
	}
	return factory
}

func (factory *rotateHookFactory) newLfsHook(options *Options) log.Hook {
	if options == nil {
		options = factory.getDefaultOption()
	}
	var (
		optArr      = options.parse()
		writer, err = rotatelogs.New(options.GetLinkName(), optArr...)
	)
	if err != nil {
		log.Errorf("config local file system for logger error: %v", err)
	}
	var (
		levels        = entity.GetLevels()
		levelEnum, ok = levels.Get(options.Level)
	)
	if !ok {
		log.SetLevel(log.WarnLevel)
	} else {
		log.SetLevel(entity.LogLevelOf(&levelEnum))
	}
	var (
		formatter = &log.TextFormatter{DisableColors: options.DisableColors}
		writerMap = lfshook.WriterMap{
			log.DebugLevel: writer, log.InfoLevel: writer, log.WarnLevel: writer,
			log.ErrorLevel: writer, log.FatalLevel: writer, log.PanicLevel: writer,
		}
		lfsHook = lfshook.NewHook(writerMap, formatter)
	)
	return lfsHook
}
