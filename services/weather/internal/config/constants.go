package config

type envKeys struct {
	ApiKey string
	Url    string
}

var EnvKeys = envKeys{
	ApiKey: "ApiKey",
	Url:    "Url",
}
