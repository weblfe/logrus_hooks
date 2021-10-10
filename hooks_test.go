package logrus_hooks

import (
	log "github.com/sirupsen/logrus"
	"github.com/weblfe/logrus_hooks/rotate"
	"github.com/weblfe/logrus_hooks/utils"
		"os"
		"testing"
	"time"
)

func TestGetMgr(t *testing.T) {
	var reg = GetMgr()
	if reg == nil {
		t.Error("get mgr impl failed")
	}
	var hook, ok = reg.Get(rotate.HookName)
	if !ok || hook == nil {
		t.Error("获取 hook factory failed")
	}
}

func TestResolve(t *testing.T) {
	var hook, ok = Resolve(rotate.HookName)
	if ok != nil || hook == nil {
		t.Error("解析构造 hook failed")
	}
}

func TestResolveAndLog(t *testing.T) {
	var (
		options = rotate.CreateOptionsWithLogName("./logs/rotate.log")
	)
	hook, ok := Resolve(rotate.HookName, options)
	if ok != nil || hook == nil {
		t.Error("解析构造 hook failed")
	}
	var (
		count  = 0
		logger = log.New()
		//	fd, err = os.OpenFile(options.LogName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
		ticker = time.NewTimer(10 * time.Second)
	)
	logger.SetOutput(os.Stdout)
	logger.AddHook(hook)
	for {
		if count >= 100 {
			ticker.Stop()
			break
		}
		select {
		case <-ticker.C:
			logger.Infoln("debug", time.Now().Format(utils.DateTimeLayout))
		default:
			logger.WithFields(log.Fields{
				"timestamp": time.Now().Unix(),
				"count":     count,
			}).Println("ticker")
			count++
		}
	}
}
