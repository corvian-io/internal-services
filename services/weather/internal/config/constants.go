package config

type envKeys struct {
	ApiKey  string
	Url     string
	ZipCode string
}

var EnvKeys = envKeys{
	ApiKey:  "ApiKey",
	Url:     "Url",
	ZipCode: "ZipCode",
}
