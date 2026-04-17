package core

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/huuhoait/gin-vue-admin/server/global"
)

// Any signing key that has ever been checked into this repository is
// considered compromised and must never authenticate a release-mode process.
var bannedJWTSigningKeys = map[string]struct{}{
	"":                                     {},
	"88c21ee8-5769-417e-b3ee-e1a309d15c88": {},
	"qmPlusv2": {},
}

// applySecretOverrides pulls secrets from the environment so that config.yaml
// can be committed with placeholders only. In release mode we refuse to start
// if any required secret is missing or still set to a known-public default.
func applySecretOverrides() {
	if key := os.Getenv("GVA_JWT_SIGNING_KEY"); key != "" {
		global.GVA_CONFIG.JWT.SigningKey = key
	}
	if pw := os.Getenv("GVA_REDIS_PASSWORD"); pw != "" {
		global.GVA_CONFIG.Redis.Password = pw
	}
	if pw := os.Getenv("GVA_DB_PASSWORD"); pw != "" {
		for i := range global.GVA_CONFIG.DBList {
			global.GVA_CONFIG.DBList[i].Password = pw
		}
		global.GVA_CONFIG.Mysql.Password = pw
		global.GVA_CONFIG.Pgsql.Password = pw
		global.GVA_CONFIG.Mssql.Password = pw
		global.GVA_CONFIG.Oracle.Password = pw
	}

	release := gin.Mode() == gin.ReleaseMode
	if _, banned := bannedJWTSigningKeys[strings.TrimSpace(global.GVA_CONFIG.JWT.SigningKey)]; banned {
		msg := "JWT signing key is empty or matches a committed default; set GVA_JWT_SIGNING_KEY"
		if release {
			panic("refusing to start: " + msg)
		}
		fmt.Println("[WARN] " + msg)
	}
}
