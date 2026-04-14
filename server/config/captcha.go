package config

type Captcha struct {
	KeyLong            int `mapstructure:"key-long" json:"key-long" yaml:"key-long"`                                     // captcha length
	ImgWidth           int `mapstructure:"img-width" json:"img-width" yaml:"img-width"`                                  // captcha width
	ImgHeight          int `mapstructure:"img-height" json:"img-height" yaml:"img-height"`                               // captcha height
	OpenCaptcha        int `mapstructure:"open-captcha" json:"open-captcha" yaml:"open-captcha"`                         // anti-brute-forcecaptchaEnableThisNumber, 0means captcha every login; other numbers mean failurespasswordTimeNumber, Such As3means captcha after three failures
	OpenCaptchaTimeOut int `mapstructure:"open-captcha-timeout" json:"open-captcha-timeout" yaml:"open-captcha-timeout"` // anti-brute-force captcha timeout, unit: s(Second)
}
