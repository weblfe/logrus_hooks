package notify

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/weblfe/logrus_hooks/entity"
	"github.com/weblfe/logrus_hooks/utils"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Options struct {
	Url         string   `json:"url" yaml:"url" env:"url"`
	Name        string   `json:"name" yaml:"name" env:"name"`
	Levels      []string `json:"level" yaml:"level" env:"level"`
	Method      string   `json:"method" yaml:"method" env:"method,post"`
	ContentType string   `json:"content_type" yaml:"content_type" env:"content_type,json"`
	logLevels   []log.Level
}

const (
	defaultHttpSchema      = "http:"
	defaultHttpMethod      = http.MethodPost
	defaultHttpContentType = ContentTypeJson
	ContentTypeJson        = "json"
	ContentTypeFrom        = "form"
	ContentTypeQuery       = "query"
	ContentTypeText        = "text"
	ContentTypeXml         = "xml"
	ContentTypePath        = "path"
)

var (
	AllSupportMethods = []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
	}
	AllSupportContentTypes = []string{
		ContentTypeJson,
		ContentTypeFrom,
		ContentTypeQuery,
		ContentTypeText,
		ContentTypeXml,
		ContentTypePath,
	}
)

func NewOptionWithEnvPrefix(prefix string, caseMode ...utils.CaseMode) *Options {
	var (
		opt = new(Options)
	)
	caseMode = append(caseMode, utils.UpperCase)
	var decoder = utils.NewEnvDecoder(caseMode[0])
	if prefix != "" {
		decoder.SetPrefix(prefix)
	}
	if err := decoder.Marshal(opt); err != nil {
		return nil
	}
	if opt.Name == "" {
		opt.Name = strings.TrimPrefix(prefix, "_")
	}
	return opt
}

func NewOptionWithEnv(caseMode utils.CaseMode, prefix ...string) *Options {
	var (
		argc    = len(prefix)
		opt     = new(Options)
		decoder = utils.NewEnvDecoder(caseMode)
	)
	if argc == 0 {
		prefix = append(prefix, "notify")
		argc++
	}
	if argc >= 0 && prefix[0] != "" {
		decoder.SetPrefix(prefix[0])
	}
	if argc >= 2 && prefix[1] != "" {
		decoder.SetPrefix(prefix[1])
	}
	if err := decoder.Marshal(opt); err != nil {
		return nil
	}
	if opt.Name == "" && prefix[0] != "" {
		opt.Name = strings.TrimPrefix(prefix[0], "_")
	}
	return opt
}

func NewOptions(arg interface{}) *Options {
	if arg == nil {
		return NewOptionWithEnvPrefix("notify")
	}
	var through = false
	for {
		switch arg.(type) {
		case []byte:
			var (
				bytes   = arg.([]byte)
				options = &Options{}
			)
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
				return NewOptionWithEnvPrefix("notify", utils.UpperCase)
			}
			// eg: case=1&prefix=app_&suffix=_logger
			if strings.Contains(key, "=") {
				var values, err = url.ParseQuery(key)
				if err != nil {
					return NewOptionWithEnvPrefix("notify", utils.UpperCase)
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
				return NewOptionWithEnv(caseMode, args...)
			}
			// eg: app_,_logger
			if strings.Contains(key, ",") {
				var kArr = strings.Split(key, ",")
				return NewOptionWithEnv(utils.UpperCase, kArr...)
			}
			// namespace prefix
			if strings.Contains(key, "_") {
				key = strings.TrimPrefix(key, "_")
			}
			NewOptionWithEnv(utils.UpperCase, key)
		case *Options:
			return arg.(*Options)
		case Options:
			var opt = arg.(Options)
			return &opt
		}
		if !through {
			break
		}
	}

	return nil
}

func (options *Options) GetLevels() []log.Level {
	// 已解析过
	if len(options.logLevels) > 0 {
		return options.logLevels
	}
	// 无限定日志level
	if len(options.Levels) <= 0 {
		options.logLevels = log.AllLevels
		return options.logLevels
	}
	var (
		enums = entity.GetLevels()
		cache = make(map[log.Level]int)
	)
	for i, v := range options.Levels {
		if enum, ok := enums.Get(v); ok {
			var level = entity.LogLevelOf(&enum)
			if _, ok := cache[level]; ok {
				continue
			}
			cache[level] = i
			options.logLevels = append(options.logLevels, level)
		}
	}
	// 指定参数异常 默认限定 warn level
	if len(options.logLevels) <= 0 {
		options.logLevels = append(options.logLevels, log.WarnLevel)
	}
	return options.logLevels
}

func (options *Options) GetUrl() string {
	if options == nil {
		return ""
	}
	if !strings.HasPrefix(options.Url, "https:") && !strings.HasPrefix(options.Url, "http:") {
		return fmt.Sprintf("%s//%s", defaultHttpSchema, options.Url)
	}
	return options.Url
}

func (options *Options) GetURI() *url.URL {
	var strUrl = options.GetUrl()
	if uri, err := url.Parse(strUrl); err == nil {
		return uri
	}
	return nil
}

func (options *Options) GetMethod() string {
	if options == nil || options.Method == "" {
		return defaultHttpMethod
	}
	var method = strings.ToLower(options.Method)
	for _, v := range AllSupportMethods {
		if v == method {
			return v
		}
	}
	return defaultHttpMethod
}

func (options *Options) GetContentType() string {
	if options == nil || options.ContentType == "" {
		return defaultHttpContentType
	}
	var contentType = strings.ToLower(options.ContentType)
	for _, v := range AllSupportContentTypes {
		if v == contentType {
			return v
		}
	}
	return defaultHttpContentType
}
