package config

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Config struct {
	CloudflareAPIToken    string `yaml:"cloudflare_api_token"`
	CloudflareAccountID   string `yaml:"cloudflare_account_id"`
	CloudflareZoneID      string `yaml:"cloudflare_zone_id"`
	CloudflareTunnelName  string `yaml:"cloudflare_tunnel_name"`
	HostName              string `yaml:"host_name"`
	AWSAccessKeyID        string `yaml:"aws_access_key_id"`
	AWSSecretAccessKey    string `yaml:"aws_secret_access_key"`
	TerraformFileLocation string `yaml:"terraform_file_location"`
	StateFileLocation     string `yaml:"state_file_location"`
}

func SetupConfig() {
	viper.SetConfigName("config")
	configDir, _ := os.UserConfigDir()
	turtleConfigDir := filepath.Join(configDir, "turtle")
	// check turtle dir exists in .config
	if _, err := os.Stat(turtleConfigDir); errors.Is(err, os.ErrNotExist) {
		// create folder
		// 0755 = rwxrwxrwx
		os.Mkdir(turtleConfigDir, 0755)
	}
	// Add search paths to find the file
	viper.AddConfigPath(turtleConfigDir)
	viper.AddConfigPath(".")
	// set type
	viper.SetConfigType("yaml")
	viper.SetConfigFile(filepath.Join(turtleConfigDir, "config.yml"))
}

func ConfigExists() (bool, error) {
	err := viper.ReadInConfig()
	// if config file isn't found
	if err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return false, nil
		} else {
			// other error
			return false, err
		}
	}
	return true, nil
}

func ReadConfig() error {
	// Find and read the config file
	err := viper.ReadInConfig()
	// if config file isn't found
	if err != nil {
		return err
	}
	return nil
}

func WriteConfig(c *Config) error {
	configBytes, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	configReader := bytes.NewReader(configBytes)
	if err = viper.ReadConfig(configReader); err != nil {
		return err
	}
	if err = viper.WriteConfig(); err != nil {
		return err
	}
	return nil
}

func GetConfigValue(key string) string {
	val := viper.Get(key)
	if val != nil {
		return val.(string)
	}
	return ""
}

func SetTFEnv() error {
	awsSecret := GetConfigValue("aws_secret_access_key")
	awsAccess := GetConfigValue("aws_access_key_id")
	if awsSecret == "" {
		return errors.New("AWS Secret Access key not set in config")
	}
	if awsAccess == "" {
		return errors.New("AWS Access key id not set in config")
	}
	os.Setenv("AWS_SECRET_ACCESS_KEY", awsSecret)
	os.Setenv("AWS_ACCESS_KEY_ID", awsAccess)
	return nil
}
