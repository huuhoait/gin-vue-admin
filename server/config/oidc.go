package config

// OIDC holds OpenID Connect provider configuration.
// Set Enabled=true and configure your provider to activate SSO.
// The application registers two routes:
//
//	GET  /base/oidc/login    — redirects browser to provider
//	GET  /base/oidc/callback — handles provider redirect, issues JWT
type OIDC struct {
	Enabled      bool   `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	ProviderURL  string `mapstructure:"provider-url" json:"provider-url" yaml:"provider-url"`
	ClientID     string `mapstructure:"client-id" json:"client-id" yaml:"client-id"`
	ClientSecret string `mapstructure:"client-secret" json:"client-secret" yaml:"client-secret"`
	RedirectURL  string `mapstructure:"redirect-url" json:"redirect-url" yaml:"redirect-url"`
	// Scopes defaults to ["openid","profile","email"] if empty
	Scopes []string `mapstructure:"scopes" json:"scopes" yaml:"scopes"`
	// ClaimUsername is the OIDC claim to use as the GVA username (default: "email")
	ClaimUsername string `mapstructure:"claim-username" json:"claim-username" yaml:"claim-username"`
}
