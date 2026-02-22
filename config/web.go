package config

// WebConfig holds the web server configuration
type WebConfig struct {
	Enable   bool   `toml:"enable" mapstructure:"enable" json:"enable"`
	Host     string `toml:"host" mapstructure:"host" json:"host"`
	Port     int    `toml:"port" mapstructure:"port" json:"port"`
	Username string `toml:"username" mapstructure:"username" json:"username"`
	Password string `toml:"password" mapstructure:"password" json:"-"`
}
