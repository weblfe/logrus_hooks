package entity

import (
		log "github.com/sirupsen/logrus"
		"testing"
)

func TestGetLevels(t *testing.T) {
		var levels =GetLevels()
		if _,ok:=levels.Get(log.DebugLevel);!ok {
				t.Error("enum get failed")
		}
}

func TestLogLevelOf(t *testing.T) {
		var levels =GetLevels()
		e,ok:=levels.Get(log.DebugLevel)
		if !ok || e.IsNull() {
				t.Error("enum get failed")
		}
		if LogLevelOf(&e) != log.DebugLevel {
				t.Error("enum level trans failed")
		}
}

