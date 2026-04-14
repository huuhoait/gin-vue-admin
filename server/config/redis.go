package config

type Redis struct {
	Name         string   `mapstructure:"name" json:"name" yaml:"name"`                         // name of the current instance
	Addr         string   `mapstructure:"addr" json:"addr" yaml:"addr"`                         // server address:port
	Password     string   `mapstructure:"password" json:"password" yaml:"password"`             // password
	DB           int      `mapstructure:"db" json:"db" yaml:"db"`                               // in single-instance moderedisofWhichPieceDatabase
	UseCluster   bool     `mapstructure:"useCluster" json:"useCluster" yaml:"useCluster"`       // YesNouse cluster mode
	ClusterAddrs []string `mapstructure:"clusterAddrs" json:"clusterAddrs" yaml:"clusterAddrs"` // node address list in cluster mode
}
