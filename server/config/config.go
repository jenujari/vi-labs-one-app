package config

// Define a struct to hold the configuration values
type Config struct {
	App struct {
		Name  string `mapstructure:"name"`
		Port  int    `mapstructure:"port"`
		Debug bool   `mapstructure:"debug"`
	} `mapstructure:"app"`
	Database struct {
		Path string `mapstructure:"path"`
	} `mapstructure:"database"`
	Secret struct {
		V1 string `mapstructure:"v1"`
	} `mapstructure:"secret"`
}
