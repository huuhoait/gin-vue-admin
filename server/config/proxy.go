package config

type Proxy struct {
	CoreServiceURL  string `mapstructure:"core_service_url" json:"core_service_url" yaml:"core_service_url"`
	OrderServiceURL string `mapstructure:"order_service_url" json:"order_service_url" yaml:"order_service_url"`
	RequestTimeout  int    `mapstructure:"request_timeout" json:"request_timeout" yaml:"request_timeout"`
	CoreDBDSN       string `mapstructure:"core_db_dsn" json:"core_db_dsn" yaml:"core_db_dsn"`
	OrderDBDSN      string `mapstructure:"order_db_dsn" json:"order_db_dsn" yaml:"order_db_dsn"`
}
