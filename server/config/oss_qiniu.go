package config

type Qiniu struct {
	Zone          string `mapstructure:"zone" json:"zone" yaml:"zone"`                                  // storage region
	Bucket        string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`                            // EmptyIntervalname
	ImgPath       string `mapstructure:"img-path" json:"img-path" yaml:"img-path"`                      // CDNacceleration domain
	AccessKey     string `mapstructure:"access-key" json:"access-key" yaml:"access-key"`                // Secret KeyAK
	SecretKey     string `mapstructure:"secret-key" json:"secret-key" yaml:"secret-key"`                // Secret KeySK
	UseHTTPS      bool   `mapstructure:"use-https" json:"use-https" yaml:"use-https"`                   // YesNoUsehttps
	UseCdnDomains bool   `mapstructure:"use-cdn-domains" json:"use-cdn-domains" yaml:"use-cdn-domains"` // UploadYesNoUseCDNupload acceleration
}
