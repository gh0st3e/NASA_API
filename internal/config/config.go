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

	NasaAPIKey = "NASA_API_KEY"

	MinioHost         = "MINIO_HOST"
	MinioPort         = "MINIO_PORT"
	MinioRootUser     = "MINIO_ROOT_USER"
	MinioRootPassword = "MINIO_ROOT_PASSWORD"
	MinioBucket       = "MINIO_BUCKET"
	MinioLocation     = "MINIO_LOCATION"
)

type Config struct {
	PSQLDatabase     PSQLDatabase
	Server           Server
	NasaClientConfig NasaClientConfig
	MinioConfig      MinioConfig
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
	APIKey string `required:"true" split_word:"true"`
}

type MinioConfig struct {
	Host       string `required:"true" split_word:"true"`
	Port       string `required:"true" split_word:"true"`
	Username   string `required:"true" split_word:"true"`
	Password   string `required:"true" split_word:"true"`
	BucketName string `required:"true" split_word:"true"`
	Location   string `required:"true" split_word:"true"`
	Address    string `required:"false"`
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

	minioConfig, err := initMinio()
	if err != nil {
		return Config{}, err
	}
	cfg.MinioConfig = minioConfig

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
		return Server{}, err
	}

	var serverConfig = Server{}

	serverConfig.Host = params[ServerHost]
	serverConfig.Port = params[ServerPort]

	serverConfig.Address = fmt.Sprintf("%s:%s", serverConfig.Host, serverConfig.Port)

	return serverConfig, nil
}

func initNasaClient() (NasaClientConfig, error) {
	var params = map[string]string{
		NasaAPIKey: "",
	}

	params, err := LookupEnvs(params)
	if err != nil {
		return NasaClientConfig{}, err
	}

	var nasaClientConfig = NasaClientConfig{}

	nasaClientConfig.APIKey = params[NasaAPIKey]

	return nasaClientConfig, nil
}

func initMinio() (MinioConfig, error) {
	var params = map[string]string{
		MinioHost:         "",
		MinioPort:         "",
		MinioRootUser:     "",
		MinioRootPassword: "",
		MinioBucket:       "",
		MinioLocation:     "",
	}

	params, err := LookupEnvs(params)
	if err != nil {
		return MinioConfig{}, err
	}

	var minioConfig = MinioConfig{}

	minioConfig.Host = params[MinioHost]
	minioConfig.Port = params[MinioPort]
	minioConfig.Username = params[MinioRootUser]
	minioConfig.Password = params[MinioRootPassword]
	minioConfig.BucketName = params[MinioBucket]
	minioConfig.Location = params[MinioLocation]
	minioConfig.Address = fmt.Sprintf("%s:%s", minioConfig.Host, minioConfig.Port)

	return minioConfig, nil
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
