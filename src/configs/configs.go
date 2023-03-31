package configs

import (
	"os"

	"github.com/joho/godotenv"
)

var Envs = GetEnvs()

func GetEnvs() map[string]string {
	envFile := os.Args[1]
	envs, _ := godotenv.Read(envFile)
	return envs
}
