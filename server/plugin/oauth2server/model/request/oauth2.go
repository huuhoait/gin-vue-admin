package request

// AuthorizeReq is the GET /oauth2/authorize query.
type AuthorizeReq struct {
	ResponseType string `form:"response_type" binding:"required,oneof=code"`
	ClientID     string `form:"client_id" binding:"required"`
	RedirectURI  string `form:"redirect_uri" binding:"required,url"`
	Scope        string `form:"scope"`
	State        string `form:"state"`
}

// TokenReq covers all three grant types — only the relevant subset of fields
// is populated per grant. Validated dynamically in the service layer.
type TokenReq struct {
	GrantType    string `form:"grant_type" json:"grant_type" binding:"required"`
	Code         string `form:"code" json:"code"`
	RedirectURI  string `form:"redirect_uri" json:"redirect_uri"`
	RefreshToken string `form:"refresh_token" json:"refresh_token"`
	Scope        string `form:"scope" json:"scope"`
	ClientID     string `form:"client_id" json:"client_id"`
	ClientSecret string `form:"client_secret" json:"client_secret"`
}

type IntrospectReq struct {
	Token string `form:"token" json:"token" binding:"required"`
}

type RevokeReq struct {
	Token string `form:"token" json:"token" binding:"required"`
}
