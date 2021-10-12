package notify

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/weblfe/logrus_hooks/entity"
	"github.com/weblfe/logrus_hooks/facede"
	"github.com/weblfe/logrus_hooks/utils"
	"time"
)

type httpHookImpl struct {
	hookUrl     string
	hookName    string
	method      string
	contentType string
	levels      []log.Level
	client      facede.WebHookClient
}

func NewHttpWebHook(options Options) *httpHookImpl {
	var hook = new(httpHookImpl)
	hook.hookName = options.Name
	hook.hookUrl = options.Url
	hook.levels = options.GetLevels()
	return hook
}

func (hook *httpHookImpl) SetClient(client facede.WebHookClient) bool {
	if hook == nil || hook.client != nil {
		return false
	}
	hook.client = client
	return true
}

func (hook *httpHookImpl) Levels() []log.Level {
	if hook == nil || hook.levels == nil {
		return log.AllLevels
	}
	return hook.levels
}

func (hook *httpHookImpl) Fire(entry *log.Entry) error {
	if hook == nil || hook.hookUrl == "" {
		return errors.New("nil hook")
	}
	if entry == nil {
		return errors.New("nil log entry")
	}
	if !hook.checkLevel(entry.Level) {
		return nil
	}
	var client = hook.resolver()
	if client == nil {
		return nil
	}
	return client.Send(hook.parseData(entry.Message,entry.Data,entry.Time))
}

func (hook *httpHookImpl) checkLevel(level log.Level) bool {
	if hook == nil || hook.levels == nil {
		return false
	}
	for _, v := range hook.levels {
		if v == level {
			return true
		}
	}
	return false
}

func (hook *httpHookImpl) GetOptions() Options {
	var opt = Options{
		Url:         hook.hookUrl,
		Method:      hook.method,
		ContentType: hook.contentType,
		Name:        hook.hookName,
		Levels:      entity.Levels(hook.levels).StringerArray(),
	}
	return opt
}

func (hook *httpHookImpl) parseData(msg string,data log.Fields,at time.Time) map[string]string {
	var kv = make(map[string]string)
	for k, v := range data {
		kv[k] = utils.NewStringer(v).String()
	}
	return kv
}

func (hook *httpHookImpl) resolver() facede.WebHookClient {
	if hook == nil {
		return nil
	}
	if hook.client == nil && hook.hookUrl != "" {
		var client = NewUrlClient(hook.hookUrl, hook.method)
		client.SetContentType(hook.contentType)
		hook.client = client
	}
	return hook.client
}
