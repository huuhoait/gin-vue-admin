package config

type System struct {
	DbType        string `mapstructure:"db-type" json:"db-type" yaml:"db-type"`    // Databasetype:mysql(default)|sqlite|sqlserver|postgresql
	OssType       string `mapstructure:"oss-type" json:"oss-type" yaml:"oss-type"` // Osstype
	RouterPrefix  string `mapstructure:"router-prefix" json:"router-prefix" yaml:"router-prefix"`
	Addr          int    `mapstructure:"addr" json:"addr" yaml:"addr"` // portValue
	LimitCountIP  int    `mapstructure:"iplimit-count" json:"iplimit-count" yaml:"iplimit-count"`
	LimitTimeIP   int    `mapstructure:"iplimit-time" json:"iplimit-time" yaml:"iplimit-time"`
	UseMultipoint bool   `mapstructure:"use-multipoint" json:"use-multipoint" yaml:"use-multipoint"`    // block concurrent logins from multiple locations
	UseRedis      bool   `mapstructure:"use-redis" json:"use-redis" yaml:"use-redis"`                   // Useredis
	UseMongo      bool   `mapstructure:"use-mongo" json:"use-mongo" yaml:"use-mongo"`                   // Usemongo
	UseStrictAuth bool   `mapstructure:"use-strict-auth" json:"use-strict-auth" yaml:"use-strict-auth"` // use treeRoleallocation mode
	DisableAutoMigrate   bool   `mapstructure:"disable-auto-migrate" json:"disable-auto-migrate" yaml:"disable-auto-migrate"`          // AutomaticmigrateMoveDatabaseTableStructure, Production EnvironmentSuggestSetForfalse, ManualmigrateMove
	// EnableSwagger registers /swagger (OpenAPI UI). When false, the UI is still
	// served in non-release Gin mode (local dev). In release mode it stays off
	// unless set to true (e.g. private staging).
	EnableSwagger bool `mapstructure:"enable-swagger" json:"enable-swagger" yaml:"enable-swagger"`
}
