package config

type GlobalConfig struct {
	Language string
	Output   string
	Format   string
	Model    string
}

var Global = &GlobalConfig{}
