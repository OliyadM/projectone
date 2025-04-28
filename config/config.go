package config

// This file is intentionally left minimal.
// All env loading is done in env.go.

// You can optionally define structs or constants here for grouped configs.

type AppConfig struct {
	DBURI     string
	DBName    string
	JWTSecret string
}

func LoadAppConfig() AppConfig {
	return AppConfig{
		DBURI:     GetEnv("MONGO_URI", "mongodb://localhost:27017"),
		DBName:    GetEnv("DB_NAME", "afro_vintage"),
		JWTSecret: GetEnv("JWT_SECRET", "fallback-secret"),
	}
}
