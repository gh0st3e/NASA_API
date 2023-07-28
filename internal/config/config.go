package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	PSQLDatabaseDriver   = "PSQL_DATABASE_DRIVER"
	PSQLDatabaseUser     = "PSQL_DATABASE_USER"
	PSQLDatabasePassword = "PSQL_DATABASE_PASSWORD"
	PSQLDatabaseHost     = "PSQL_DATABASE_HOST"
	PSQLDatabasePort     = "PSQL_DATABASE_PORT"
	PSQLDatabaseName     = "PSQL_DATABASE_NAME"

	ServerHost = "SERVER_HOST"
	ServerPort = "SERVER_PORT"

	NasaApiKey = "NASA_API_KEY"
)

type Config struct {
	PSQLDatabase     PSQLDatabase
	Server           Server
	NasaClientConfig NasaClientConfig
}

type PSQLDatabase struct {
	Driver       string `required:"true" split_word:"true"`
	User         string `required:"true" split_word:"true"`
	Password     string `required:"true" split_word:"true"`
	Host         string `required:"true" split_word:"true"`
	Port         string `required:"true" split_word:"true"`
	Name         string `required:"true" split_word:"true"`
	Timeout      int    `required:"true" split_word:"true"`
	DefaultLimit int    `required:"true" split_word:"true"`
	DefaultPage  int    `required:"true" split_word:"true"`
	Address      string `required:"false"`
}

type Server struct {
	Host    string `required:"true" split_word:"true"`
	Port    string `required:"true" split_word:"true"`
	Address string `required:"false"`
}

type NasaClientConfig struct {
	ApiKey string `required:"true" split_word:"true"`
}

func Init() (Config, error) {
	// .env - for docker
	// local.env -  for local load
	err := godotenv.Load(".env")
	if err != nil {
		return Config{}, err
	}

	var cfg = Config{}

	psql, err := initPSQL()
	if err != nil {
		return Config{}, err
	}
	cfg.PSQLDatabase = psql

	serverConfig, err := initServer()
	if err != nil {
		return Config{}, err
	}
	cfg.Server = serverConfig

	nasaClientConfig, err := initNasaClient()
	if err != nil {
		return Config{}, err
	}
	cfg.NasaClientConfig = nasaClientConfig

	return cfg, nil
}

func initPSQL() (PSQLDatabase, error) {
	var params = map[string]string{
		PSQLDatabaseDriver:   "",
		PSQLDatabaseUser:     "",
		PSQLDatabasePassword: "",
		PSQLDatabaseHost:     "",
		PSQLDatabasePort:     "",
		PSQLDatabaseName:     "",
	}

	params, err := LookupEnvs(params)
	if err != nil {
		return PSQLDatabase{}, err
	}
	var db = PSQLDatabase{}

	db.Driver = params[PSQLDatabaseDriver]
	db.User = params[PSQLDatabaseUser]
	db.Password = params[PSQLDatabasePassword]
	db.Host = params[PSQLDatabaseHost]
	db.Port = params[PSQLDatabasePort]
	db.Name = params[PSQLDatabaseName]

	db.Address = fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", db.Driver, db.User, db.Password, db.Host, db.Port, db.Name)

	return db, nil
}

func initServer() (Server, error) {
	var params = map[string]string{
		ServerHost: "",
		ServerPort: "",
	}

	params, err := LookupEnvs(params)
	if err != nil {
		return Server{}, nil
	}

	var serverConfig = Server{}

	serverConfig.Host = params[ServerHost]
	serverConfig.Port = params[ServerPort]

	serverConfig.Address = fmt.Sprintf("%s:%s", serverConfig.Host, serverConfig.Port)

	return serverConfig, nil
}

func initNasaClient() (NasaClientConfig, error) {
	var params = map[string]string{
		NasaApiKey: "",
	}

	params, err := LookupEnvs(params)
	if err != nil {
		return NasaClientConfig{}, nil
	}

	var nasaClientConfig = NasaClientConfig{}

	nasaClientConfig.ApiKey = params[NasaApiKey]

	return nasaClientConfig, nil
}

func LookupEnvs(params map[string]string) (map[string]string, error) {
	var errorMsg string

	for i := range params {
		envVar, ok := os.LookupEnv(i)
		if !ok {
			errorMsg += fmt.Sprintf("\nCannot find %s", i)
		}
		params[i] = envVar
	}

	if len(errorMsg) > 0 {
		return nil, fmt.Errorf("%s", errorMsg)
	}
	return params, nil
}
