package configs

type PgSQLConnectionParams struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

type ServerConfig struct {
	PortToStart string                `toml:"start_port"`
	ConnParams  PgSQLConnectionParams `toml:"server"`
}

func CreateConfigForServer() *ServerConfig {
	return &ServerConfig{}
}
