package env

import (
	"github.com/joho/godotenv"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	DatabaseConnection = os.Getenv("DB_CONN")
	ServeAddr = os.Getenv("SERVE_ADDR")
}

var (
	DatabaseConnection string
	ServeAddr          string
)
