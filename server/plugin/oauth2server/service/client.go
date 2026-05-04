package service

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/plugin/oauth2server/model"
	"github.com/huuhoait/gin-vue-admin/server/plugin/oauth2server/model/request"
	"github.com/huuhoait/gin-vue-admin/server/utils"

	"gorm.io/datatypes"
)

type clientService struct{}

// CreateClient registers a new OAuth2 client. Returns the plaintext secret
// alongside the persisted record — the secret is shown to the admin once and
// never readable afterwards.
func (s *clientService) CreateClient(req request.CreateClientReq, createdBy uint) (model.OAuth2Client, string, error) {
	clientID, err := secureToken(16)
	if err != nil {
		return model.OAuth2Client{}, "", err
	}
	secret, err := secureToken(32)
	if err != nil {
		return model.OAuth2Client{}, "", err
	}
	uris, _ := json.Marshal(req.RedirectURIs)
	grants, _ := json.Marshal(s.normalizeGrants(req.GrantTypes))
	scopes, _ := json.Marshal(req.Scopes)
	row := model.OAuth2Client{
		ClientID:         clientID,
		ClientSecretHash: utils.BcryptHash(secret),
		Name:             req.Name,
		Description:      req.Description,
		RedirectURIs:     datatypes.JSON(uris),
		GrantTypes:       datatypes.JSON(grants),
		Scopes:           datatypes.JSON(scopes),
		Enabled:          true,
		CreatedBy:        createdBy,
	}
	if err := global.GVA_DB.Create(&row).Error; err != nil {
		return row, "", err
	}
	return row, secret, nil
}

func (s *clientService) UpdateClient(req request.UpdateClientReq) (model.OAuth2Client, error) {
	var row model.OAuth2Client
	if err := global.GVA_DB.Where("id = ?", req.ID).First(&row).Error; err != nil {
		return row, err
	}
	updates := map[string]any{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.RedirectURIs != nil {
		uris, _ := json.Marshal(req.RedirectURIs)
		updates["redirect_uris"] = datatypes.JSON(uris)
	}
	if req.GrantTypes != nil {
		grants, _ := json.Marshal(s.normalizeGrants(req.GrantTypes))
		updates["grant_types"] = datatypes.JSON(grants)
	}
	if req.Scopes != nil {
		scopes, _ := json.Marshal(req.Scopes)
		updates["scopes"] = datatypes.JSON(scopes)
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}
	if len(updates) == 0 {
		return row, nil
	}
	if err := global.GVA_DB.Model(&row).Updates(updates).Error; err != nil {
		return row, err
	}
	return row, global.GVA_DB.Where("id = ?", req.ID).First(&row).Error
}

func (s *clientService) DeleteClient(id uint) error {
	return global.GVA_DB.Delete(&model.OAuth2Client{}, "id = ?", id).Error
}

func (s *clientService) FindByID(id uint) (model.OAuth2Client, error) {
	var c model.OAuth2Client
	err := global.GVA_DB.Where("id = ?", id).First(&c).Error
	return c, err
}

func (s *clientService) FindByClientID(clientID string) (model.OAuth2Client, error) {
	var c model.OAuth2Client
	err := global.GVA_DB.Where("client_id = ?", clientID).First(&c).Error
	return c, err
}

func (s *clientService) ListClients(req request.ClientListReq) ([]model.OAuth2Client, int64, error) {
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	db := global.GVA_DB.Model(&model.OAuth2Client{})
	if kw := strings.TrimSpace(req.Keyword); kw != "" {
		db = db.Where("name LIKE ? OR client_id LIKE ?", "%"+kw+"%", "%"+kw+"%")
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []model.OAuth2Client
	err := db.Limit(req.PageSize).Offset(req.PageSize * (req.Page - 1)).Order("id DESC").Find(&list).Error
	return list, total, err
}

// RegenerateSecret rotates the client secret. Existing tokens stay valid;
// the new secret only affects future /oauth2/token calls.
func (s *clientService) RegenerateSecret(id uint) (string, error) {
	secret, err := secureToken(32)
	if err != nil {
		return "", err
	}
	if err := global.GVA_DB.Model(&model.OAuth2Client{}).Where("id = ?", id).
		Update("client_secret_hash", utils.BcryptHash(secret)).Error; err != nil {
		return "", err
	}
	return secret, nil
}

// AuthenticateClient verifies clientID + plaintext secret. Returns the client
// row when valid; ErrClientAuth otherwise.
var ErrClientAuth = errors.New("invalid client credentials")

func (s *clientService) AuthenticateClient(clientID, secret string) (model.OAuth2Client, error) {
	c, err := s.FindByClientID(clientID)
	if err != nil {
		return c, ErrClientAuth
	}
	if !c.Enabled {
		return c, ErrClientAuth
	}
	if !utils.BcryptCheck(secret, c.ClientSecretHash) {
		return c, ErrClientAuth
	}
	return c, nil
}

// ValidateRedirect checks the URI against the client's whitelist using exact
// string match (no wildcard or partial matches per RFC 6749 §3.1.2.2).
func (s *clientService) ValidateRedirect(c model.OAuth2Client, uri string) bool {
	var uris []string
	if err := json.Unmarshal(c.RedirectURIs, &uris); err != nil {
		return false
	}
	return contains(uris, uri)
}

func (s *clientService) supportsGrant(c model.OAuth2Client, grant string) bool {
	var grants []string
	if err := json.Unmarshal(c.GrantTypes, &grants); err != nil {
		return false
	}
	return contains(grants, grant)
}

// normalizeGrants drops empty entries and dedups so the column reflects a
// canonical set.
func (s *clientService) normalizeGrants(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, g := range in {
		g = strings.TrimSpace(g)
		if g == "" {
			continue
		}
		if _, dup := seen[g]; dup {
			continue
		}
		seen[g] = struct{}{}
		out = append(out, g)
	}
	return out
}

