package configs

import "os"

type Env struct {
	Port     string
	MongoURI string
	DBName   string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func NewEnv() *Env {
	return &Env{
		Port:     getEnv("PORT", "8000"),
		MongoURI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DBName:   getEnv("DB_NAME", "goagra"),
	}
}
