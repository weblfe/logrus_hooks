package notify

type Options struct {
		Url string `json:"url" yaml:"url" env:"url"`
		Name string `json:"name" yaml:"name" env:"name"`
		Levels []string `json:"level" yaml:"level" env:"level"`
}

