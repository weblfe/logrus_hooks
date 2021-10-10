package faced

// WebHookClient 客户端
type WebHookClient interface {
		Send(params map[string]string) error
}
