package config

// Define a struct to hold the configuration values
type Config struct {
	App struct {
		Name  string `mapstructure:"name"`
		Port  int    `mapstructure:"port"`
		Debug bool   `mapstructure:"debug"`
	} `mapstructure:"app"`
	Secret struct {
		Secret string
		// POSTGRES_PASSWORD string
		// POSTGRES_DB       string
		// POSTGRES_USER     string
	}
}
