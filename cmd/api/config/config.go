package config

type DbConfig struct {
	Host            string
	User            string
	Password        string
	Database        string
	ApplicationName string
}

type KafkaProducerConfig struct {
	BootstrapServers string
	Acks             string
}

type HttpServerConfig struct {
	Addr string
}

type AppConfig struct {
	DbConfig         DbConfig
	KPConfig         KafkaProducerConfig
	HttpServerConfig HttpServerConfig
}
