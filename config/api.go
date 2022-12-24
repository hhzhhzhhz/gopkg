package config

var defaultConfiguration = New()

func LoadFromDataSource(content []byte, unmarshal Unmarshaller) error {
	return defaultConfiguration.reflush(content, unmarshal)
}

func GetConfig() *Config {
	return defaultConfiguration
}
