package config

type ContextKey string

const (
	REPO_KEY ContextKey = "CONTEXT_REPOSITORY"
)

// Define a struct to hold the configuration values
type Config struct {
	App struct {
		Name  string `mapstructure:"name"`
		Port  int    `mapstructure:"port"`
		Debug bool   `mapstructure:"debug"`
	} `mapstructure:"app"`
	Database struct {
		Domain string `mapstructure:"domain"`
		Port   int    `mapstructure:"port"`
	} `mapstructure:"database"`
	Secret struct {
		UserName          string
		Password          string
		Secret            string
		ApiKey            string
		ApiSecret         string
		POSTGRES_PASSWORD string
		POSTGRES_DB       string
		POSTGRES_USER     string
	}
}
