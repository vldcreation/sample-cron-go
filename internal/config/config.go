package config

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App     AppConfig     `mapstructure:"app" yaml:"app,omitempty"`
	Minio   MinioConfig   `mapstructure:"minio" yaml:"minio,omitempty"`
	GCS     GCSConfig     `mapstructure:"gcs" yaml:"gcs,omitempty"`
	Storage StorageConfig `mapstructure:"storage" yaml:"storage,omitempty"`
	PubSub  PubSubConfig  `mapstructure:"pubsub" yaml:"pubsub,omitempty"`
}

func NewAppConfig() *Config {
	// env := ENV()

	// require := []string{
	// 	"APP_ENV",
	// 	"APP_PORT",
	// 	"APP_HOST",
	// 	"APP_NAME",
	// }

	// if err := env.Require(require...); err != nil {
	// 	log.Fatalln("missing environment variables", require)
	// 	<-time.After(time.Second * 5)
	// 	panic(err)
	// }

	viper.AddConfigPath("./")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("error reading config file", err)
		<-time.After(time.Second * 5)
		panic(err)
	}

	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalln("error unmarshalling config", err)
		<-time.After(time.Second * 5)
		panic(err)
	}

	// in this case, i only need to set the env from the first & second argument to implement switcher
	// @TODO: implement dynamic load env from args
	if len(os.Args) > 2 {
		os.Setenv("APP_ENV", os.Args[1])
		os.Setenv("STORAGE_BUCKET", os.Args[2])
		config.App.APP_ENV = os.Getenv("APP_ENV")
		config.Storage.Bucket = os.Getenv("STORAGE_BUCKET")
	}

	return &config

}

type AppConfig struct {
	APP_ENV  string `mapstructure:"app_env" yaml:"app_env,omitempty"`
	APP_PORT string `mapstructure:"app_port" yaml:"app_port,omitempty"`
	APP_HOST string `mapstructure:"app_host" yaml:"app_host,omitempty"`
	APP_NAME string `mapstructure:"app_name" yaml:"app_name,omitempty"`
}

type MinioConfig struct {
	Endpoint  string `mapstructure:"minio_endpoint" yaml:"minio_endpoint,omitempty"`
	AccessKey string `mapstructure:"minio_access_key" yaml:"minio_access_key,omitempty"`
	SecretKey string `mapstructure:"minio_secret_key" yaml:"minio_secret_key,omitempty"`
	Bucket    string `mapstructure:"minio_bucket" yaml:"minio_bucket,omitempty"`
	Prefix    string `mapstructure:"minio_prefix" yaml:"minio_prefix,omitempty"`
	Region    string `mapstructure:"minio_region" yaml:"minio_region,omitempty"`
	UseSSL    bool   `mapstructure:"minio_use_ssl" yaml:"minio_use_ssl,omitempty"`
}

type GCSConfig struct {
	AccountPath string `mapstructure:"gcs_account_path" yaml:"gcs_account_path" json:"gcs_account_path"`
	Bucket      string `mapstructure:"gcs_bucket" yaml:"gcs_bucket" json:"gcs_bucket"`
	Prefix      string `mapstructure:"gcs_prefix" yaml:"gcs_prefix" json:"gcs_prefix"`
	AcecssID    string `mapstructure:"gcs_access_id" yaml:"gcs_access_id" json:"gcs_access_id"`
	PrivateKey  string `mapstructure:"gcs_private_key" yaml:"gcs_private_key" json:"gcs_private_key"`
}

// PubSubConfig
type PubSubConfig struct {
	AccountPath  string `mapstructure:"pubsub_account_path" yaml:"pubsub_account_path" json:"pubsub_account_path"`
	ProjectID    string `mapstructure:"pubsub_project_id" yaml:"pubsub_project_id" json:"pubsub_project_id"`
	Topic        string `mapstructure:"pubsub_topic" yaml:"pubsub_topic" json:"pubsub_topic"`
	Subscription string `mapstructure:"pubsub_subscription" yaml:"pubsub_subscription" json:"pubsub_subscription"`
}

// GeneralConfig fields for storage for switcher purpose
type StorageConfig struct {
	Bucket string `yaml:"storage_bucket" json:"storage_bucket"`
	Prefix string `yaml:"storage_prefix" json:"storage_prefix"`
}
