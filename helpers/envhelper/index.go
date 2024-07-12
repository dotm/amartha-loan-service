package envhelper

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func GetEnvVar(key string) string {
	osEnv := os.Getenv(key)
	if osEnv != "" {
		return osEnv
	}
	return viper.GetString(key)
}

func SetLocalEnvVar() (err error) {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("error reading config file: %v\n", err)
		os.Exit(1)
	}

	return nil
}
