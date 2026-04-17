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

	enforceDBTLS(release)
}

// enforceDBTLS ensures the database connection is encrypted in transit. In
// release mode we refuse to start if the config looks plaintext; in debug we
// only warn so local setups keep working.
//
// The checks are heuristic: we look at the driver-specific Config/DSN hint
// that GORM appends to the connection string.
//   - Postgres: must contain `sslmode=verify-full`, `verify-ca`, or `require`
//   - MySQL:    must contain `tls=true` or `tls=preferred` (or a named TLS
//     config registered by the operator)
func enforceDBTLS(release bool) {
	pgCfg := strings.ToLower(global.GVA_CONFIG.Pgsql.Config)
	if global.GVA_CONFIG.Pgsql.Path != "" {
		ok := strings.Contains(pgCfg, "sslmode=require") ||
			strings.Contains(pgCfg, "sslmode=verify-ca") ||
			strings.Contains(pgCfg, "sslmode=verify-full")
		if !ok {
			msg := "postgres connection is not TLS-protected (config missing sslmode=require|verify-ca|verify-full)"
			if release {
				panic("refusing to start: " + msg)
			}
			fmt.Println("[WARN] " + msg)
		}
	}
	myCfg := strings.ToLower(global.GVA_CONFIG.Mysql.Config)
	if global.GVA_CONFIG.Mysql.Path != "" {
		// tls=false or tls=skip-verify is explicitly rejected; anything
		// else (tls=true, tls=preferred, tls=<custom>) is accepted.
		disabled := strings.Contains(myCfg, "tls=false") ||
			strings.Contains(myCfg, "tls=skip-verify") ||
			!strings.Contains(myCfg, "tls=")
		if disabled {
			msg := "mysql connection is not TLS-protected (config missing tls=true or tls=<named>)"
			if release {
				panic("refusing to start: " + msg)
			}
			fmt.Println("[WARN] " + msg)
		}
	}
}
