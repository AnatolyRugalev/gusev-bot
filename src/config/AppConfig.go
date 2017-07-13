package config

type MongoConfig struct {
	Host     string `json:"host"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type BotConfig struct {
	Token string `json:"token"`
}

type AppConfig struct {
	Bot BotConfig `json:"bot"`
	Mongo MongoConfig `json:"mongo"`
}
