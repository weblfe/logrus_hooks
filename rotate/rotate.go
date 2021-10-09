package rotate

import (
	"encoding/json"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"github.com/weblfe/logrus_hooks/entity"
	"time"
)

type (
	rotateHookFactory struct {
		defaultOption *Options
	}

	Options struct {
		RotationCount uint          `json:"rotation_count" yaml:"rotate_count" xml:"rotate_count"`
		DisableColors bool          `json:"disable_colors" yaml:"disable_colors" xml:"disable_colors"`
		LogNameLayout string        `json:"log_name_layout" yaml:"log_name_layout" xml:"log_name_layout"`
		LogName       string        `json:"log_name" yaml:"log_name" xml:"log_name"`
		RotationTime  time.Duration `json:"rotation_time" yaml:"rotation_time" xml:"rotation_time"`
		MaxAge        time.Duration `json:"max_age" yaml:"max_age" xml:"max_age"`
		Level         string        `json:"level" yaml:"level" xml:"level"`
		optArr        []rotatelogs.Option
	}
)

func CreateRotateFactory() *rotateHookFactory {
	var factory = new(rotateHookFactory)
	return factory
}

func newOption(data []byte) *Options {
	var options = new(Options)
	if json.Valid(data) {
		if err := json.Unmarshal(data, options); err != nil {
			return nil
		}
	}
	return options
}

func (factory *rotateHookFactory) Create(args ...[]byte) log.Hook {
	if len(args) == 0 {
		return factory.newLfsHook(factory.getDefaultOption())
	}
	var options = newOption(args[0])
	return factory.newLfsHook(options)
}

func (factory *rotateHookFactory) getDefaultOption() *Options {
	if factory.defaultOption == nil {
		factory.defaultOption = &Options{}
	}
	return factory.defaultOption
}

func (factory *rotateHookFactory) newLfsHook(options *Options) log.Hook {
	if options == nil {
		options = factory.getDefaultOption()
	}
	var (
		optArr      = options.parse()
		writer, err = rotatelogs.New(options.LogNameLayout, optArr...)
	)
	if err != nil {
		log.Errorf("config local file system for logger error: %v", err)
	}
	var (
		levelEnum, ok = entity.GetLevels().Get(options.Level)
	)
	if ok {
		var level = entity.LogLevelOf(&levelEnum)
		log.SetLevel(level)
	} else {
		log.SetLevel(log.WarnLevel)
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

func (option *Options) parse() []rotatelogs.Option {
	if option.optArr != nil {
		return option.optArr
	}
	var optArr []rotatelogs.Option
	// WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
	optArr = append(optArr, rotatelogs.WithLinkName(option.LogName))
	// WithRotationTime设置日志分割的时间，这里设置为一小时分割一次
	if option.RotationTime > 0 {
		optArr = append(optArr, rotatelogs.WithRotationTime(option.RotationTime))
	}
	// WithMaxAge和WithRotationCount二者只能设置一个，
	// WithMaxAge设置文件清理前的最长保存时间，
	// WithRotationCount设置文件清理前最多保存的个数。
	if option.MaxAge > 0 {
		optArr = append(optArr, rotatelogs.WithMaxAge(option.MaxAge))
	}
	if option.RotationCount > 0 && option.MaxAge <= 0 {
		optArr = append(optArr, rotatelogs.WithRotationCount(option.RotationCount))
	}
	if optArr != nil && len(optArr) > 0 {
		option.optArr = optArr
	}
	return optArr
}
