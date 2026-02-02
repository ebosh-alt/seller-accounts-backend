package config

type Config struct {
	Server   ServerConfig   `yaml:"Server" mapstructure:"Server"`
	Postgres PostgresConfig `yaml:"Postgres" mapstructure:"Postgres"`
	Cache    CacheConfig    `yaml:"Cache" mapstructure:"Cache"`
}

type PostgresConfig struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     string `yaml:"port" mapstructure:"port"`
	User     string `yaml:"user" mapstructure:"user"`
	Password string `yaml:"password" mapstructure:"password"`
	DBName   string `yaml:"DBName" mapstructure:"DBName"`
	SSLMode  string `yaml:"sslMode" mapstructure:"sslMode"`
	PgDriver string `yaml:"pgDriver" mapstructure:"pgDriver"`
}

type ServerConfig struct {
	AppVersion string    `yaml:"appVersion" mapstructure:"appVersion"`
	Host       string    `yaml:"host" mapstructure:"host" validate:"required"`
	Port       string    `yaml:"port" mapstructure:"port" validate:"required"`
	TLS        TLSConfig `yaml:"tls" mapstructure:"tls"`
}

type CacheConfig struct {
	TTLSeconds int `yaml:"ttlSeconds" mapstructure:"ttlSeconds"`
}

type TLSConfig struct {
	Enabled    bool             `yaml:"enabled" mapstructure:"enabled"`
	CertFile   string           `yaml:"certFile" mapstructure:"certFile"`
	KeyFile    string           `yaml:"keyFile" mapstructure:"keyFile"`
	SelfSigned SelfSignedConfig `yaml:"selfSigned" mapstructure:"selfSigned"`
}

type SelfSignedConfig struct {
	Enabled      bool     `yaml:"enabled" mapstructure:"enabled"`
	CommonName   string   `yaml:"commonName" mapstructure:"commonName"`
	Hosts        []string `yaml:"hosts" mapstructure:"hosts"`
	ValidForDays int      `yaml:"validForDays" mapstructure:"validForDays"`
}
