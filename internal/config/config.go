package config

import "os"

type Config struct {
	DBUrl string
}

func Load() Config {
	dbUrl := os.Getenv("DB_URL")

	if dbUrl == "" {
		dbUrl = "file:controlplane.db?_pragma=busy_timeout(5000)"
	}

	return Config{
		DBUrl: dbUrl,
	}
}
