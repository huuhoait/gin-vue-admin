package system

import (
	"context"

	"go.uber.org/zap"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/system"
	"github.com/huuhoait/gin-vue-admin/server/utils"
)

type JwtService struct{}

var JwtServiceApp = new(JwtService)

// JsonInBlacklist persists the revoked token to the database (audit log) and
// publishes it to the distributed blacklist so every instance sees the
// revocation immediately. Redis is the source of truth at request time; the
// DB row is only used to rebuild Redis if it is flushed.
func (jwtService *JwtService) JsonInBlacklist(jwtList system.JwtBlacklist) (err error) {
	err = global.GVA_DB.Create(&jwtList).Error
	if err != nil {
		return
	}
	return utils.BlacklistAdd(jwtList.Jwt)
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetRedisJWT
//@description: retrieve JWT from Redis
//@param: userName string
//@return: redisJWT string, err error

func (jwtService *JwtService) GetRedisJWT(userName string) (redisJWT string, err error) {
	redisJWT, err = global.GVA_REDIS.Get(context.Background(), userName).Result()
	return redisJWT, err
}

// LoadAll rehydrates the blacklist on startup. Tokens that have already
// expired are skipped — BlacklistAdd uses each token's remaining TTL, so
// expired ones would either set a negative TTL or burn keys in Redis forever.
func LoadAll() {
	var data []string
	err := global.GVA_DB.Model(&system.JwtBlacklist{}).Select("jwt").Find(&data).Error
	if err != nil {
		global.GVA_LOG.Error("Failed to load JWT blacklist from database!", zap.Error(err))
		return
	}
	for _, tok := range data {
		if utils.BlacklistTTL(tok) <= 0 {
			continue
		}
		_ = utils.BlacklistAdd(tok)
	}
}
