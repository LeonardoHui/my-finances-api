package configs

import (
	"os"

	"github.com/joho/godotenv"
)

var Envs = GetEnvs()

func GetEnvs() map[string]string {
	if len(os.Args) > 1 {
		envFile := os.Args[1]
		envs, _ := godotenv.Read(envFile)
		return envs
	}
	return nil
}
