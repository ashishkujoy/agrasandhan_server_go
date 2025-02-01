package configs

import "os"

type Env struct {
	Port          string
	DBName        string
	MongoURI      string
	RedisURI      string
	SessionSecret string
	JWTKey        []byte
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func NewEnv() *Env {
	return &Env{
		Port:          getEnv("PORT", "8000"),
		DBName:        getEnv("DB_NAME", "goagra"),
		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		RedisURI:      getEnv("REDIS_URI", "localhost:6379"),
		SessionSecret: getEnv("SESSION_SECRET", "secret"),
		JWTKey:        []byte(getEnv("SESSION_KEY", "jwt_key")),
	}
}
