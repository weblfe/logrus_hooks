package rotate

import (
	"encoding/json"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"github.com/weblfe/logrus_hooks/utils"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type (
	// Options 构建hook 参数
	Options struct {
		RotationCount uint          `json:"rotation_count" yaml:"rotate_count" env:"rotate_count,20"`
		DisableColors bool          `json:"disable_colors" yaml:"disable_colors" env:"disable_colors,false"`
		LogNameLayout string        `json:"log_name_layout" yaml:"log_name_layout" env:"log_name_layout,%s-%Y%m%d.log"`
		LogName       string        `json:"log_name" yaml:"log_name" env:"log_name,app"`
		RotationTime  time.Duration `json:"rotation_time" yaml:"rotation_time" env:"rotation_time,24h"`
		MaxAge        time.Duration `json:"max_age" yaml:"max_age" env:"max_age,0"`
		Level         string        `json:"level" yaml:"level" env:"level,warn"`
		optArr        []rotatelogs.Option
	}

)

func NewOption(data ...interface{}) *Options {
	var options = new(Options)
	if len(data) <= 0 {
		return CreateOptionsWithEnv(utils.UpperCase)
	}
	var (
		arg     = data[0]
		through = false
	)
	for {
		switch arg.(type) {
		case []byte:
			var bytes = arg.([]byte)
			if json.Valid(bytes) {
				if err := json.Unmarshal(bytes, options); err == nil {
					return options
				}
			}
			arg = string(bytes)
			through = true
		case string:
			through = false
			var key = arg.(string)
			if key == "" {
				return CreateOptionsWithEnv(utils.UpperCase)
			}
			// eg: case=1&prefix=app_&suffix=_logger
			if strings.Contains(key, "=") {
				var values, err = url.ParseQuery(key)
				if err != nil {
					return CreateOptionsWithEnv(utils.UpperCase)
				}
				var (
					prefix, suffix string
					caseMode       = utils.UnDefineCase
					args           = make([]string, 2)
				)
				for k, v := range values {
					if len(v) <= 0 {
						continue
					}
					var value = v[0]
					switch strings.ToLower(k) {
					case "case":
						if n, e := strconv.Atoi(value); e == nil {
							caseMode = utils.CaseMode(n)
							continue
						}
						switch strings.ToLower(value) {
						case "upper":
							caseMode = utils.UpperCase
						case "lower":
							caseMode = utils.LowerCase
						case "normal":
							caseMode = utils.NormalCase
						case "default":
							caseMode = utils.UpperCase
						}
					case "prefix":
						prefix = strings.TrimSpace(value)
					case "suffix":
						suffix = strings.TrimSpace(value)
					}
				}
				if caseMode == utils.UnDefineCase {
					caseMode = utils.UpperCase
				}
				if prefix != "" {
					args[0] = prefix
				}
				if suffix != "" {
					args[1] = suffix
				}
				return CreateOptionsWithEnv(caseMode, args...)
			}
			// eg: app_,_logger
			if strings.Contains(key, ",") {
				var kArr = strings.Split(key, ",")
				return CreateOptionsWithEnv(utils.UpperCase, kArr...)
			}
			// namespace prefix
			if !strings.Contains(key, "_") {
				key = key + "_"
			}
			CreateOptionsWithEnv(utils.UpperCase, key)
		case *Options:
			return arg.(*Options)
		default:
			break
		}
		if !through {
			break
		}
	}
	return CreateOptionsWithEnv(utils.UpperCase)
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

func (option *Options) GetLinkName() string {
	var (
		name   = option.LogName
		layout = option.LogNameLayout
	)
	if strings.Contains(name, LogExt) {
		name = strings.TrimSuffix(name, LogExt)
	}
	if option.LogNameLayout == "" {
		return name + "-%Y%m%d" + LogExt
	}
	if strings.HasPrefix(option.LogNameLayout, "%s") {
		layout = strings.Replace(option.LogNameLayout, "%s", name, 1)
	}
	if layout != "" && !strings.HasSuffix(layout, LogExt) {
		layout = layout + LogExt
	}
	return layout
}

func (option *Options) String() string {
	var bytes = option.Bytes()
	if len(bytes) <= 0 {
		return ""
	}
	return string(bytes)
}

func (option *Options) Bytes() []byte {
	if bytes, err := json.Marshal(option); err == nil {
		return bytes
	}
	return nil
}

func CreateDefaultOptions() *Options {
	var opt = new(Options)
	opt.LogName = "app.log"
	opt.MaxAge = 0
	opt.Level = "warn"
	opt.DisableColors = false
	opt.RotationCount = 20
	opt.LogNameLayout = `%s-%Y%m%d.log`
	opt.RotationTime = 24 * time.Hour
	return opt
}

func CreateOptionsWithLogName(name string) *Options {
	var opt = new(Options)
	opt.MaxAge = 0
	opt.LogName = name
	opt.Level = "warn"
	opt.DisableColors = false
	opt.RotationCount = 20
	opt.LogNameLayout = `%s-%Y%m%d` + LogExt
	opt.RotationTime = 24 * time.Hour
	return opt
}

func CreateOptionsWithEnv(caseMode utils.CaseMode, prefix ...string) *Options {
	var (
		argc    = len(prefix)
		opt     = new(Options)
		decoder = utils.NewEnvDecoder(caseMode)
	)
	if argc >= 1 && prefix[0] != "" {
		decoder.SetPrefix(prefix[0])
	}
	if argc >= 2 && prefix[1] != "" {
		decoder.SetSuffix(prefix[1])
	}
	if err := decoder.Marshal(opt); err != nil {
		log.Error("decoder option error:", err)
	}
	return opt
}
