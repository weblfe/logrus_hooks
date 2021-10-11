package facede

import log "github.com/sirupsen/logrus"

// WebHookClient 客户端
type WebHookClient interface {
	Send(params map[string]string) error
}

type NotifierHook interface {
	log.Hook
	SetClient(client WebHookClient) bool
}
