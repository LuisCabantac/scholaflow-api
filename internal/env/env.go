package env

import (
	"log"

	goenv "github.com/allisson/go-env"
)

func GetString(key, fallback string) string {
	return goenv.GetString(key, fallback)
}

func GetRequired(key string) string {
	val := goenv.GetString(key, "")
	if val == "" {
		log.Fatalf("missing required environment variable: %s", key)
	}

	return val
}
