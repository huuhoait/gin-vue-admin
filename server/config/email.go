package config

type Email struct {
	To          string `mapstructure:"to" json:"to" yaml:"to"`                               // recipient:multiple values separated by English commas Example:a@qq.com b@qq.com use this project as a parameter in real development
	From        string `mapstructure:"from" json:"from" yaml:"from"`                         // sender  the mailbox you send mail from
	Host        string `mapstructure:"host" json:"host" yaml:"host"`                         // server address For example smtp.qq.com  please go toQQor the mailbox you send fromviewItssmtpProtocol
	Secret      string `mapstructure:"secret" json:"secret" yaml:"secret"`                   // secret key    used for loginsecret key avoid using the mailboxpassword go to the mailbox providersmtpapply one for login used assecret key
	Nickname    string `mapstructure:"nickname" json:"nickname" yaml:"nickname"`             // nickname    sendernickname usually your own mailbox
	Port        int    `mapstructure:"port" json:"port" yaml:"port"`                         // port     please go toQQor the mailbox you send fromviewItssmtpProtocol mostly 465
	IsSSL       bool   `mapstructure:"is-ssl" json:"is-ssl" yaml:"is-ssl"`                   // enable SSL   YesNoEnableSSL
	IsLoginAuth bool   `mapstructure:"is-loginauth" json:"is-loginauth" yaml:"is-loginauth"` // enable LoginAuth   YesNoUseLoginAuthAuthenticationMethod(SuitableUsed forIBM, MicrosoftEmailServeretc.)
}
