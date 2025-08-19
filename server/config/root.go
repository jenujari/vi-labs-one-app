package config

import (
	"log"
	"os"
	"path/filepath"

	// "github.com/goforj/godump"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	DB_NAME = "data.db"
)


var (
	dbPATH  string
	dirPATH string
	dbc     *gorm.DB
	logger  *log.Logger
	config  *Config
)

func init() {
	viper.SetConfigName("conf") // Name of the config file (without extension)
	viper.SetConfigType("yaml") // Config file type
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Read the config file
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Initialize variables by unmarshaling into the struct
	config = new(Config)
	err = viper.Unmarshal(config)
	if err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}

	config.Secret.UserName = os.Getenv("ZERODHA_USER")
	config.Secret.Password = os.Getenv("ZERODHA_PASS")
	config.Secret.Secret = os.Getenv("ZERODHA_SECRET")
	config.Secret.ApiKey = os.Getenv("ZERODHA_API_KEY")
	config.Secret.ApiSecret = os.Getenv("ZERODHA_API_SECRET")

	// init log system
	logger = log.Default()
	logger.SetOutput(os.Stdout)

	// set default config folder in user default directory

	dirPATH = filepath.Join(config.Database.Path)
	dbPATH = filepath.Join(dirPATH, DB_NAME)
	os.MkdirAll(dirPATH, os.ModePerm)

	// setup db connection
	dbc, err = gorm.Open(sqlite.Open(dbPATH), &gorm.Config{})
	if err != nil {
		logger.Panic("failed to connect database on path " + dbPATH + ": " + err.Error())
	}

	// godump.Dump(config)
}

func GetDBC() *gorm.DB {
	return dbc
}

func GetUserDir() string {
	return dirPATH
}

func GetLogger() *log.Logger {
	return logger
}

func GetConfig() *Config {
	return config
}
