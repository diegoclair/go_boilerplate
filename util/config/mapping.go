package config

type Config struct {
	App AppConfig `mapstructure:"app"`
	DB  DBConfig  `mapstructure:"db"`
	Log LogConfig `mapstructure:"log"`
}

type AppConfig struct {
	Name        string     `mapstructure:"name"`
	Environment string     `mapstructure:"environment"`
	Port        string     `mapstructure:"port"`
	Auth        AuthConfig `mapstructure:"auth"`
}
type AuthConfig struct {
	JWTPrivateKey string `mapstructure:"jwt-private-key"`
}

type DBConfig struct {
	MySQL MySQLConfig `mapstructure:"mysql"`
}

type MySQLConfig struct {
	Username           string `mapstructure:"username"`
	Password           string `mapstructure:"password"`
	Host               string `mapstructure:"host"`
	Port               string `mapstructure:"port"`
	DBName             string `mapstructure:"db-name"`
	CryptoKey          string `mapstructure:"crypto-key"`
	MaxLifeInMinutes   int    `mapstructure:"max-life-in-minutes"`
	MaxIdleConnections int    `mapstructure:"max-idle-connections"`
	MaxOpenConnections int    `mapstructure:"max-open-connections"`
}

type LogConfig struct {
	Debug     bool   `mapstructure:"debug"`
	LogToFile bool   `mapstructure:"log-to-file"`
	Path      string `mapstructure:"path"`
}
