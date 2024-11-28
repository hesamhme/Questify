package config

// Config represents the application configuration.
type Config struct {
	DB  DatabaseConfig `yaml:"db"`
	JWT JWTConfig      `yaml:"jwt"`
}

// DatabaseConfig holds database-related configuration.
type DatabaseConfig struct {
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	DBName string `yaml:"db_name"`
}

// JWTConfig holds JWT-related configuration.
type JWTConfig struct {
	Secret string `yaml:"secret"`
}
