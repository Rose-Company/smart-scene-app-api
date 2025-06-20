package config

import (
	"fmt"
	"os"
	"strings"

	_ "embed"

	l "smart-scene-app-api/pkg/logger"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds all settings
//
//go:embed config.yaml
var defaultConfig []byte

type Schema struct {
	Postgres struct {
		Host      string `mapstructure:"host"`
		Port      string `mapstructure:"port"`
		User      string `mapstructure:"user"`
		Pass      string `mapstructure:"pass"`
		Db        string `mapstructure:"db"`
		Params    string `mapstructure:"params"`
		GormDebug string `mapstructure:"gorm_debug"`
	} `mapstructure:"postgres"`

	Redis *Redis `yaml:"redis" mapstructure:"redis"`

	Google struct {
		CredentialsDir string `mapstructure:"credentials_dir"`
	} `mapstructure:"google"`

	DigitalOcean struct {
		StorageAccessKey     string `mapstructure:"storage_access_key"`
		StorageSecretKey     string `mapstructure:"storage_secret_key"`
		StorageEndPoint      string `mapstructure:"storage_endpoint"`
		StorageRegion        string `mapstructure:"storage_region"`
		StorageBucket        string `mapstructure:"storage_bucket"`
		ImgkitOutputEndpoint string `mapstructure:"imgkit_output_endpoint"`
		StorageAcl           string `mapstructure:"storage_acl"`
	} `mapstructure:"digital_ocean"`
	Http struct {
		MaxIdleConnection     int `mapstructure:"max_idle_connection"`
		IdleConnectionTimeout int `mapstructure:"idle_connection_timeout"`
	} `mapstructure:"http"`

	AwsSes struct {
		Region          string `mapstructure:"region"`
		AccessKeyID     string `mapstructure:"access_key_id"`
		SecretAccessKey string `mapstructure:"secret_access_key"`
	} `mapstructure:"aws_ses"`

	JwtSecret        string `mapstructure:"jwt_secret"`
	TokenExpiredTime int64  `mapstructure:"token_expired_time"`

	JWT struct {
		SecretKey string `mapstructure:"secret_key"`
	} `mapstructure:"jwt"`
}

// Redis ...
type Redis struct {
	Host string `yaml:"host" mapstructure:"host"`
	Port string `yaml:"internal_port" mapstructure:"internal_port"`
	DB   int    `yaml:"db_idx" mapstructure:"db_idx"`
	Pass string `yaml:"pass" mapstructure:"pass"`
}

var Config Schema

func init() {
	ll := l.New()

	// ✅ Load .env file trước
	if err := godotenv.Load(); err != nil {
		ll.Info("No .env file found, continuing...")
	}

	// ✅ Expand biến môi trường trong config.yaml (nhúng sẵn)
	expanded := os.ExpandEnv(string(defaultConfig))
	
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(strings.NewReader(expanded))
	if err != nil {
		ll.Fatal("Failed to read viper config", l.Error(err))
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()

	err = viper.Unmarshal(&Config)
	if err != nil {
		ll.Fatal("Failed to unmarshal config", l.Error(err))
	}
}

func GetPartKey(partID int) string {
	return fmt.Sprintf("part_%v", partID)
}
