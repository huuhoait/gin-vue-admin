package service

import (
	"context"

	"github.com/huuhoait/gin-vue-admin/server/global"
	sysModel "github.com/huuhoait/gin-vue-admin/server/model/system"
)

// ResolveUserNames maps GVA admin user UUIDs (as strings) to a display name.
// Missing UUIDs are simply absent from the result map. Errors are swallowed —
// enrichment is best-effort and must never block the proxied response.
//
// NickName is preferred; falls back to Username. An empty/zero UUID is skipped
// so placeholder values don't hit the DB.
func ResolveUserNames(ctx context.Context, uuids []string) map[string]string {
	out := map[string]string{}
	if len(uuids) == 0 || global.GVA_DB == nil {
		return out
	}

	seen := make(map[string]struct{}, len(uuids))
	filtered := make([]string, 0, len(uuids))
	for _, u := range uuids {
		if u == "" || u == "00000000-0000-0000-0000-000000000000" {
			continue
		}
		if _, dup := seen[u]; dup {
			continue
		}
		seen[u] = struct{}{}
		filtered = append(filtered, u)
	}
	if len(filtered) == 0 {
		return out
	}

	var rows []sysModel.SysUser
	if err := global.GVA_DB.WithContext(ctx).
		Select("uuid", "username", "nick_name").
		Where("uuid IN ?", filtered).
		Find(&rows).Error; err != nil {
		return out
	}

	for _, r := range rows {
		name := r.NickName
		if name == "" {
			name = r.Username
		}
		out[r.UUID.String()] = name
	}
	return out
}
