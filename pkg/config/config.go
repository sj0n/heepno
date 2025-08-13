package config

type GlobalConfig struct {
	Language      string
	Output        string
	Format        string
	AaiModel      string
	DeepgramModel string
	OpenaiModel   string
}

var Global = &GlobalConfig{}
