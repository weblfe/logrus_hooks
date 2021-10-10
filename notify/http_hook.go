package notify

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/weblfe/logrus_hooks/entity"
	"github.com/weblfe/logrus_hooks/faced"
	"github.com/weblfe/logrus_hooks/utils"
)

type httpHookImpl struct {
	hookUrl     string
	hookName    string
	method      string
	contentType string
	level       []string
	cacheLevels map[log.Level]int
	client      faced.WebHookClient
}

func NewHttpWebHook(options Options) *httpHookImpl {
	var hook = new(httpHookImpl)
	hook.hookName = options.Name
	hook.hookUrl = options.Url
	hook.level = options.Levels
	hook.cacheLevels = make(map[log.Level]int)
	return hook
}

func (hook *httpHookImpl) SetClient(client faced.WebHookClient) bool {
	if hook == nil || hook.client != nil {
		return false
	}
	hook.client = client
	return true
}

func (hook *httpHookImpl) Levels() []log.Level {
	if hook == nil || hook.level == nil {
		return log.AllLevels
	}
	var levels []log.Level
	if len(hook.cacheLevels) > 0 {
		for level := range hook.cacheLevels {
			levels = append(levels, level)
		}
		return levels
	}
	var levelMgr = entity.GetLevels()
	for i, v := range hook.level {
		var enum, ok = levelMgr.Get(v)
		if !ok {
			continue
		}
		var level = entity.LogLevelOf(&enum)
		if _, ok = hook.cacheLevels[level]; ok {
			continue
		}
		levels = append(levels, level)
		hook.cacheLevels[level] = i
	}
	return levels
}

func (hook *httpHookImpl) Fire(entry *log.Entry) error {
	if hook == nil || hook.hookUrl == "" {
		return errors.New("nil hook")
	}
	if entry == nil {
		return errors.New("nil log entry")
	}
	if _, ok := hook.cacheLevels[entry.Level]; !ok {
		return nil
	}
	var client = hook.resolver()
	if client == nil {
		return nil
	}
	return client.Send(hook.parseData(entry.Data))
}

func (hook *httpHookImpl) parseData(data log.Fields) map[string]string {
	var kv = make(map[string]string)
	for k, v := range data {
		kv[k] = utils.NewStringer(v).String()
	}
	return kv
}


func (hook *httpHookImpl) resolver() faced.WebHookClient {
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
