package entity

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

const (
	LogLevel Symbol = "level"
)

var (
	levelEnums = NewEnumMgr(LogLevel)
)

type (
	Levels []log.Level
)

func init() {
	levelEnums.load(initLogLevels)
}

// 注册消息
func initLogLevels(mgr *EnumMgr) {
	mgr.Add(NewEnum("panic", LogLevel, "FatalLevel").SetCustom(log.PanicLevel))
	mgr.Add(NewEnum("fatal", LogLevel, "FatalLevel").SetCustom(log.FatalLevel))
	mgr.Add(NewEnum("error", LogLevel, "ErrorLevel").SetCustom(log.ErrorLevel))
	mgr.Add(NewEnum("warn", LogLevel, "WarnLevel").SetCustom(log.WarnLevel))
	mgr.Add(NewEnum("info", LogLevel, "InfoLevel").SetCustom(log.InfoLevel))
	mgr.Add(NewEnum("debug", LogLevel, "DebugLevel").SetCustom(log.DebugLevel))
	mgr.Add(NewEnum("trace", LogLevel, "TraceLevel").SetCustom(log.TraceLevel))
}

func GetLevels() *EnumMgr {
	return levelEnums
}

func LogLevelOf(e *Enum) log.Level {
	if e == nil {
		return log.WarnLevel
	}
	if e.symbol != LogLevel {
		return log.WarnLevel
	}
	var v = e.GetCustom()
	switch v.(type) {
	case log.Level:
		return v.(log.Level)
	case uint32:
		return log.Level(v.(uint32))
	}
	return log.WarnLevel
}

func (levels Levels) StringerArray() []string {
	var (
		result []string
		enums  = GetLevels()
	)
	for _, v := range levels {
		if enum, ok := enums.Get(v); ok {
			result = append(result, fmt.Sprintf("%v", enum.value))
		}
	}
	return result
}
