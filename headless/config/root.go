package config

import (
	"log"
	"os"

	// "github.com/goforj/godump"

	// "github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

var (
	// dbc          *pgxpool.Pool
	logger *log.Logger
	config *Config
	// PGSQL_STRING = "postgres://%s:%s@%s:%d/%s"
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

	config.Secret.Secret = os.Getenv("ZERODHA_SECRET")

	// init log system
	logger = log.Default()
	logger.SetOutput(os.Stdout)

	// godump.Dump(config)
}

// func SetuDbConnection(ctx context.Context) {
// 	var err error

// 	PGSQL_STRING = fmt.Sprintf(PGSQL_STRING, config.Secret.POSTGRES_USER, config.Secret.POSTGRES_PASSWORD, config.Database.Domain, config.Database.Port, config.Secret.POSTGRES_DB)

// 	poolConfig, err := pgxpool.ParseConfig(PGSQL_STRING)
// 	if err != nil {
// 		panic("Unable to parse database config: " + err.Error())
// 	}

// 	poolConfig.MaxConns = 10
// 	poolConfig.MinConns = 2

// 	dbc, err = pgxpool.NewWithConfig(ctx, poolConfig)
// 	if err != nil {
// 		panic("Unable to connect to database: " + err.Error())
// 	}

// 	err = dbc.Ping(ctx)
// 	if err != nil {
// 		panic("Unable to ping database: " + err.Error())
// 	}

// 	fmt.Println("Connected to database successfully!")
// }

// func CloseDbConnection() {
// 	if dbc != nil {
// 		dbc.Close()
// 	}
// }

// func GetDBC() *pgxpool.Pool {
// 	return dbc
// }

func GetLogger() *log.Logger {
	return logger
}

func GetConfig() *Config {
	return config
}
