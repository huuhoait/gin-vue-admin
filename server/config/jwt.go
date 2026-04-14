package config

type JWT struct {
	SigningKey  string `mapstructure:"signing-key" json:"signing-key" yaml:"signing-key"`    // jwtSign
	ExpiresTime string `mapstructure:"expires-time" json:"expires-time" yaml:"expires-time"` // expiration time
	BufferTime  string `mapstructure:"buffer-time" json:"buffer-time" yaml:"buffer-time"`    // buffer time
	Issuer      string `mapstructure:"issuer" json:"issuer" yaml:"issuer"`                   // issuer
}
