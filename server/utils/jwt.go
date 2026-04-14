package utils

import (
	"context"
	"errors"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	jwt "github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	SigningKey []byte
}

var (
	TokenValid            = errors.New("UnknownError")
	TokenExpired          = errors.New("tokenAlreadyExpired")
	TokenNotValidYet      = errors.New("tokenStillNotActivate")
	TokenMalformed        = errors.New("ThisNotYesOnePiecetoken")
	TokenSignatureInvalid = errors.New("NoneEffectSign")
	TokenInvalid          = errors.New("UnableHandleThistoken")
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.GVA_CONFIG.JWT.SigningKey),
	}
}

func (j *JWT) CreateClaims(baseClaims request.BaseClaims) request.CustomClaims {
	bf, _ := ParseDuration(global.GVA_CONFIG.JWT.BufferTime)
	ep, _ := ParseDuration(global.GVA_CONFIG.JWT.ExpiresTime)
	claims := request.CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: int64(bf / time.Second), // buffer time1Day buffer timeInsideWillGetNewoftokenRefresh Token ThisWhenOnePieceUserWillExistsTwoPieceHaveEffectToken ButYesFrontendOnlykeepOnePiece otherOnePieceWillLost
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"GVA"},                   // Receivepublic
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)), // SignEffective Time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ep)),    // expiration time 7Day  configurationFile
			Issuer:    global.GVA_CONFIG.JWT.Issuer,              // SignofsendRowPerson
		},
	}
	return claims
}

// CreateToken CreateOnePiecetoken
func (j *JWT) CreateToken(claims request.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// CreateTokenByOldToken Oldtoken ChangeNewtoken UsereturnAndReturnSourceAvoidConcurrencyIssue
func (j *JWT) CreateTokenByOldToken(oldToken string, claims request.CustomClaims) (string, error) {
	v, err, _ := global.GVA_Concurrency_Control.Do("JWT:"+oldToken, func() (interface{}, error) {
		return j.CreateToken(claims)
	})
	return v.(string), err
}

// ParseToken Parse token
func (j *JWT) ParseToken(tokenString string) (*request.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &request.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, TokenExpired
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, TokenMalformed
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, TokenSignatureInvalid
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, TokenNotValidYet
		default:
			return nil, TokenInvalid
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*request.CustomClaims); ok && token.Valid {
			return claims, nil
		}
	}
	return nil, TokenValid
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetRedisJWT
//@description: jwtSaveInredisAndsetexpiration time
//@param: jwt string, userName string
//@return: err error

func SetRedisJWT(jwt string, userName string) (err error) {
	// This Locationexpiration timeetc.Atjwtexpiration time
	dr, err := ParseDuration(global.GVA_CONFIG.JWT.ExpiresTime)
	if err != nil {
		return err
	}
	timer := dr
	err = global.GVA_REDIS.Set(context.Background(), userName, jwt, timer).Err()
	return err
}
