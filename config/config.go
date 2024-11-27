package config

type Config struct {
	DB DatabaseConfig `yaml:"db"`
}

type DatabaseConfig struct {
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	DBName string `yaml:"db_name"`
}
