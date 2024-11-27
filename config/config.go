package config

type Config struct {
	DB DatabaseConfig `yaml:"db"`
	JWT JWTConfig      `yaml:"jwt"`
}

type DatabaseConfig struct {
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	DBName string `yaml:"db_name"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
}