package config

import (
	"strings"

	"gorm.io/gorm/logger"
)

type DsnProvider interface {
	Dsn() string
}

// Embeded StructureBodyCanByflattenToUpperOneLayer, FromAndprotectHold config FileofStructureAndOriginalComeOneSample
// See playground: https://go.dev/play/p/KIcuhqEoxmY

// GeneralDB AlsoBy Pgsql And Mysql OriginalSampleUse
type GeneralDB struct {
	Prefix       string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                         // database table prefix
	Port         string `mapstructure:"port" json:"port" yaml:"port"`                               // Databaseport
	Config       string `mapstructure:"config" json:"config" yaml:"config"`                         // advanced configuration
	Dbname       string `mapstructure:"db-name" json:"db-name" yaml:"db-name"`                      // DatabaseName
	Username     string `mapstructure:"username" json:"username" yaml:"username"`                   // database username
	Password     string `mapstructure:"password" json:"password" yaml:"password"`                   // Databasepassword
	Path         string `mapstructure:"path" json:"path" yaml:"path"`                               // database host
	Engine       string `mapstructure:"engine" json:"engine" yaml:"engine" default:"InnoDB"`        // database engine; defaultInnoDB
	LogMode      string `mapstructure:"log-mode" json:"log-mode" yaml:"log-mode"`                   // whether to enable global GORM SQL log
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"max-idle-conns" yaml:"max-idle-conns"` // max idle connections
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"max-open-conns" yaml:"max-open-conns"` // max open connections to the database
	Singular     bool   `mapstructure:"singular" json:"singular" yaml:"singular"`                   // whether to disable plural table names globally; true enables
	LogZap       bool   `mapstructure:"log-zap" json:"log-zap" yaml:"log-zap"`                      // whether to write logs to file via zap
}

func (c GeneralDB) LogLevel() logger.LogLevel {
	switch strings.ToLower(c.LogMode) {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Info
	}
}

type SpecializedDB struct {
	Type      string `mapstructure:"type" json:"type" yaml:"type"`
	AliasName string `mapstructure:"alias-name" json:"alias-name" yaml:"alias-name"`
	GeneralDB `yaml:",inline" mapstructure:",squash"`
	Disable   bool `mapstructure:"disable" json:"disable" yaml:"disable"`
}
