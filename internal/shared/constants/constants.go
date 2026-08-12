package constants

var EnvKeys = envKeys{
	Env:        "ENV",
	DBHost:     "DB_HOST",
	DBPort:     "DB_PORT",
	DBUser:     "DB_USER",
	DBPassword: "DB_PASSWORD",
	DBName:     "DB_NAME",
	DBSSLMode:  "DB_SSLMODE",
	DBSchema:   "DB_SCHEMA",
}

type envKeys struct {
	Env        string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	DBSchema   string
}
