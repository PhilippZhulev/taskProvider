package configurator

// Config ...
// Протокол конфигураций
type Config struct {
	ServiceADDR    string `toml:"service_addr"`
	DatabaseURL    string `toml:"database_url"`
	Salt           string `toml:"salt"`
}

// NewConfig ...
// Задать конфигурации
func NewConfig() *Config {
	return &Config{}
}
