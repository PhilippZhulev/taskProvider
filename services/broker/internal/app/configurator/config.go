package configurator

// Config ...
// Протокол конфигураций
type Config struct {
	ServiceADDR    string `toml:"service_addr"`
	UserADDR    string `toml:"user_addr"`
}

// NewConfig ...
// Задать конфигурации
func NewConfig() *Config {
	return &Config{}
}
