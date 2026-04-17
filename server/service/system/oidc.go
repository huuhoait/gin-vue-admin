package system

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	gooidc "github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"gorm.io/gorm"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/system"
	"github.com/huuhoait/gin-vue-admin/server/utils"
)

// OIDCService handles OpenID Connect authentication flows.
type OIDCService struct{}

// oidcProvider and oidcConfig are initialised lazily by GetOIDCConfig.
var (
	oidcProvider *gooidc.Provider
	oidcConfig   *oauth2.Config
)

// GetOIDCConfig returns the OAuth2 config for the configured provider.
// Returns an error if OIDC is disabled or misconfigured.
func (s *OIDCService) GetOIDCConfig(ctx context.Context) (*oauth2.Config, *gooidc.Provider, error) {
	cfg := global.GVA_CONFIG.OIDC
	if !cfg.Enabled {
		return nil, nil, errors.New("OIDC is not enabled")
	}
	if oidcProvider == nil {
		var err error
		oidcProvider, err = gooidc.NewProvider(ctx, cfg.ProviderURL)
		if err != nil {
			return nil, nil, fmt.Errorf("OIDC provider init: %w", err)
		}
		scopes := cfg.Scopes
		if len(scopes) == 0 {
			scopes = []string{gooidc.ScopeOpenID, "profile", "email"}
		}
		oidcConfig = &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Endpoint:     oidcProvider.Endpoint(),
			Scopes:       scopes,
		}
	}
	return oidcConfig, oidcProvider, nil
}

// LoginURL returns the provider redirect URL with a state parameter.
func (s *OIDCService) LoginURL(ctx context.Context, state string) (string, error) {
	oauth2cfg, _, err := s.GetOIDCConfig(ctx)
	if err != nil {
		return "", err
	}
	return oauth2cfg.AuthCodeURL(state, oauth2.AccessTypeOnline), nil
}

// HandleCallback exchanges the authorization code for tokens, verifies the ID
// token, and returns (or auto-creates) the corresponding SysUser.
func (s *OIDCService) HandleCallback(ctx context.Context, code string) (*system.SysUser, error) {
	oauth2cfg, provider, err := s.GetOIDCConfig(ctx)
	if err != nil {
		return nil, err
	}

	token, err := oauth2cfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange: %w", err)
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token in response")
	}

	verifier := provider.Verifier(&gooidc.Config{ClientID: oauth2cfg.ClientID})
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("id_token verification: %w", err)
	}

	var claims struct {
		Sub   string `json:"sub"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("claims extraction: %w", err)
	}

	providerName := global.GVA_CONFIG.OIDC.ProviderURL
	claimUser := global.GVA_CONFIG.OIDC.ClaimUsername
	if claimUser == "" {
		claimUser = "email"
	}
	username := claims.Email
	if claimUser == "sub" {
		username = claims.Sub
	}

	// Find existing user by OAuthSub
	var user system.SysUser
	err = global.GVA_DB.WithContext(ctx).
		Where("oauth_provider = ? AND oauth_sub = ?", providerName, claims.Sub).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Auto-provision: create a new user linked to this OIDC identity.
		user = system.SysUser{
			UUID:          uuid.New(),
			Username:      username,
			NickName:      claims.Name,
			Email:         claims.Email,
			Password:      utils.BcryptHash(uuid.New().String()), // unusable random password
			AuthorityId:   uint(888),
			OAuthProvider: providerName,
			OAuthSub:      claims.Sub,
			Enable:        1,
		}
		if err := global.GVA_DB.WithContext(ctx).Create(&user).Error; err != nil {
			return nil, fmt.Errorf("auto-provision user: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("db lookup: %w", err)
	}

	// Preload authority for JWT claims
	if err := global.GVA_DB.WithContext(ctx).
		Preload("Authorities").Preload("Authority").
		Where("id = ?", user.ID).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
