package notify

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/weblfe/logrus_hooks/entity"
	"github.com/weblfe/logrus_hooks/utils"
	"net/http"
	"net/url"
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
	AllSupportContentTypes = []string{
		ContentTypeJson,
		ContentTypeFrom,
		ContentTypeQuery,
		ContentTypeText,
		ContentTypeXml,
		ContentTypePath,
	}

	AllSupportMethods = []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
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
	return opt
}

func NewOptions(arg interface{}) *Options {

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
