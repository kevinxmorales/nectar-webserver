package env

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

func LoadEnvironment() error {
	currentEnv := os.Getenv("ENV")
	if currentEnv == "" {
		currentEnv = "develop"
	}
	log.Infof("using %s environment variables", currentEnv)
	return godotenv.Load(currentEnv + ".env")
}
