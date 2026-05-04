package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/plugin/oauth2server/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	authCodeTTL    = 5 * time.Minute
	accessTTL      = 1 * time.Hour
	refreshTTL     = 30 * 24 * time.Hour
	tokenTypeBear  = "Bearer"
)

var (
	ErrInvalidGrant       = errors.New("invalid_grant")
	ErrInvalidRequest     = errors.New("invalid_request")
	ErrUnsupportedGrant   = errors.New("unsupported_grant_type")
	ErrInvalidScope       = errors.New("invalid_scope")
	ErrInvalidRedirectURI = errors.New("invalid redirect_uri for this client")
	ErrClientDisabled     = errors.New("client disabled")
)

type oauth2Service struct{}

// IssueCode is called by GET /oauth2/authorize after the admin's JWT has
// been validated. It does NOT do consent prompting — auto-approves under the
// assumption that an admin granting access to their own clients is the only
// flow this server needs to support today. Add a consent UI before exposing
// this to end-user accounts.
func (s *oauth2Service) IssueCode(client model.OAuth2Client, userID uint, redirectURI, scope string) (string, error) {
	if !client.Enabled {
		return "", ErrClientDisabled
	}
	if !Service.Client.ValidateRedirect(client, redirectURI) {
		return "", ErrInvalidRedirectURI
	}
	if !Service.Client.supportsGrant(client, "authorization_code") {
		return "", ErrUnsupportedGrant
	}
	if err := s.checkScope(client, scope); err != nil {
		return "", err
	}
	code, err := secureToken(32)
	if err != nil {
		return "", err
	}
	row := model.OAuth2AuthCode{
		Code:        code,
		ClientID:    client.ClientID,
		UserID:      userID,
		RedirectURI: redirectURI,
		Scope:       scope,
		ExpiresAt:   time.Now().Add(authCodeTTL).UTC(),
	}
	if err := global.GVA_DB.Create(&row).Error; err != nil {
		return "", err
	}
	return code, nil
}

// TokenResponse is the JSON body of a successful /oauth2/token call.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// ExchangeCode handles grant_type=authorization_code. Marks the code used in
// the same transaction as token creation so the code is single-use even
// under concurrent retries.
func (s *oauth2Service) ExchangeCode(client model.OAuth2Client, code, redirectURI string) (TokenResponse, error) {
	var resp TokenResponse
	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var ac model.OAuth2AuthCode
		// FOR UPDATE serialises concurrent redemptions of the same code at the
		// DB layer. Without it READ COMMITTED lets two transactions both see
		// used_at IS NULL and proceed — RFC 6749 §10.5 requires single-use.
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("code = ?", code).First(&ac).Error; err != nil {
			return ErrInvalidGrant
		}
		if ac.UsedAt != nil {
			return ErrInvalidGrant
		}
		if time.Now().After(ac.ExpiresAt) {
			return ErrInvalidGrant
		}
		if ac.ClientID != client.ClientID {
			return ErrInvalidGrant
		}
		if ac.RedirectURI != redirectURI {
			return ErrInvalidGrant
		}
		now := time.Now().UTC()
		if err := tx.Model(&ac).Update("used_at", now).Error; err != nil {
			return err
		}
		issued, err := s.issueToken(tx, client.ClientID, ac.UserID, ac.Scope, "authorization_code", true)
		if err != nil {
			return err
		}
		resp = issued
		return nil
	})
	return resp, err
}

func (s *oauth2Service) ClientCredentials(client model.OAuth2Client, scope string) (TokenResponse, error) {
	if !Service.Client.supportsGrant(client, "client_credentials") {
		return TokenResponse{}, ErrUnsupportedGrant
	}
	if err := s.checkScope(client, scope); err != nil {
		return TokenResponse{}, err
	}
	var resp TokenResponse
	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		issued, err := s.issueToken(tx, client.ClientID, 0, scope, "client_credentials", false)
		if err != nil {
			return err
		}
		resp = issued
		return nil
	})
	return resp, err
}

// RefreshAccessToken implements grant_type=refresh_token with token rotation:
// the presented refresh token is revoked and a new pair is issued. Mitigates
// refresh-token replay if a leaked token is later used.
func (s *oauth2Service) RefreshAccessToken(client model.OAuth2Client, refreshToken string) (TokenResponse, error) {
	if !Service.Client.supportsGrant(client, "refresh_token") {
		return TokenResponse{}, ErrUnsupportedGrant
	}
	var resp TokenResponse
	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var tok model.OAuth2Token
		// Lock the row so concurrent refresh attempts can't both observe
		// revoked_at IS NULL and mint two replacement pairs.
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("refresh_token = ?", refreshToken).First(&tok).Error; err != nil {
			return ErrInvalidGrant
		}
		if tok.RevokedAt != nil {
			return ErrInvalidGrant
		}
		if tok.RefreshExpiresAt != nil && time.Now().After(*tok.RefreshExpiresAt) {
			return ErrInvalidGrant
		}
		if tok.ClientID != client.ClientID {
			return ErrInvalidGrant
		}
		now := time.Now().UTC()
		if err := tx.Model(&tok).Update("revoked_at", now).Error; err != nil {
			return err
		}
		issued, err := s.issueToken(tx, client.ClientID, tok.UserID, tok.Scope, "refresh_token", true)
		if err != nil {
			return err
		}
		resp = issued
		return nil
	})
	return resp, err
}

// IntrospectionResult mirrors RFC 7662 §2.2 (the canonical fields). Inactive
// tokens return only {active:false}.
type IntrospectionResult struct {
	Active    bool   `json:"active"`
	ClientID  string `json:"client_id,omitempty"`
	UserID    uint   `json:"user_id,omitempty"`
	Scope     string `json:"scope,omitempty"`
	TokenType string `json:"token_type,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	GrantType string `json:"grant_type,omitempty"`
}

func (s *oauth2Service) Introspect(_ context.Context, accessOrRefresh string) IntrospectionResult {
	var tok model.OAuth2Token
	err := global.GVA_DB.Where("access_token = ? OR refresh_token = ?", accessOrRefresh, accessOrRefresh).First(&tok).Error
	if err != nil {
		return IntrospectionResult{Active: false}
	}
	if tok.RevokedAt != nil {
		return IntrospectionResult{Active: false}
	}
	// If the lookup matched on access_token, check its expiry. On
	// refresh_token, check the refresh expiry.
	now := time.Now()
	if tok.AccessToken == accessOrRefresh && now.After(tok.AccessExpiresAt) {
		return IntrospectionResult{Active: false}
	}
	if tok.RefreshToken == accessOrRefresh && tok.RefreshExpiresAt != nil && now.After(*tok.RefreshExpiresAt) {
		return IntrospectionResult{Active: false}
	}
	return IntrospectionResult{
		Active:    true,
		ClientID:  tok.ClientID,
		UserID:    tok.UserID,
		Scope:     tok.Scope,
		TokenType: tokenTypeBear,
		ExpiresAt: tok.AccessExpiresAt.Unix(),
		GrantType: tok.GrantType,
	}
}

// Revoke marks any matching access OR refresh token revoked. Idempotent —
// unknown tokens still return nil per RFC 7009 §2.2.
func (s *oauth2Service) Revoke(token string) error {
	now := time.Now().UTC()
	return global.GVA_DB.Model(&model.OAuth2Token{}).
		Where("(access_token = ? OR refresh_token = ?) AND revoked_at IS NULL", token, token).
		Update("revoked_at", now).Error
}

// PruneExpired drops tokens whose refresh window has elapsed. Auth codes
// auto-expire and are pruned at the same time. Called from the cron task.
func (s *oauth2Service) PruneExpired() (int64, error) {
	now := time.Now()
	res := global.GVA_DB.Where("expires_at < ?", now).Delete(&model.OAuth2AuthCode{})
	if res.Error != nil {
		return 0, res.Error
	}
	tokRes := global.GVA_DB.Where("(refresh_expires_at IS NOT NULL AND refresh_expires_at < ?) OR (refresh_expires_at IS NULL AND access_expires_at < ?)", now, now).
		Delete(&model.OAuth2Token{})
	return res.RowsAffected + tokRes.RowsAffected, tokRes.Error
}

func (s *oauth2Service) issueToken(tx *gorm.DB, clientID string, userID uint, scope, grant string, withRefresh bool) (TokenResponse, error) {
	access, err := secureToken(32)
	if err != nil {
		return TokenResponse{}, err
	}
	now := time.Now().UTC()
	row := model.OAuth2Token{
		AccessToken:     access,
		ClientID:        clientID,
		UserID:          userID,
		Scope:           scope,
		GrantType:       grant,
		AccessExpiresAt: now.Add(accessTTL),
	}
	resp := TokenResponse{
		AccessToken: access,
		TokenType:   tokenTypeBear,
		ExpiresIn:   int64(accessTTL.Seconds()),
		Scope:       scope,
	}
	if withRefresh {
		refresh, err := secureToken(32)
		if err != nil {
			return TokenResponse{}, err
		}
		exp := now.Add(refreshTTL)
		row.RefreshToken = refresh
		row.RefreshExpiresAt = &exp
		resp.RefreshToken = refresh
	}
	if err := tx.Create(&row).Error; err != nil {
		return TokenResponse{}, err
	}
	return resp, nil
}

func (s *oauth2Service) checkScope(client model.OAuth2Client, requested string) error {
	if requested == "" {
		return nil
	}
	var allowed []string
	if err := json.Unmarshal(client.Scopes, &allowed); err != nil {
		return ErrInvalidScope
	}
	if len(allowed) == 0 {
		return nil // client has no scope restrictions
	}
	for _, want := range strings.Fields(requested) {
		if !contains(allowed, want) {
			return ErrInvalidScope
		}
	}
	return nil
}
