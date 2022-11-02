package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claim struct {
	jwt.StandardClaims
	TokenId   string `json:"jti"`
	Tokentype string `json:"token_type"`
	UserId    int    `json:"user_id"`
}

type JWTToken struct {
	Token    string
	Claim    Claim
	ExpireAt time.Time
	Scheme   string
}

type JWT interface {
	CreateAccessToken(claim Claim) (*JWTToken, error)
	CreateRefreshToken(claim Claim) (*JWTToken, error)
}

type JWTConfig struct {
	TokenAccessExpiration  time.Duration
	TokenRefreshExpiration time.Duration
}

const (
	secretKey    = "secret"
	bearerScheme = "bearer"
	accessToken  = "access"
	refreshToken = "refresh"
)

func NewJWT(config JWTConfig) *JWTConfig {
	return &JWTConfig{
		TokenAccessExpiration:  config.TokenAccessExpiration,
		TokenRefreshExpiration: config.TokenRefreshExpiration,
	}
}

func (j *JWTConfig) CreateAccessToken(claim Claim) (*JWTToken, error) {
	exp := time.Now().Add(j.TokenAccessExpiration * time.Minute)
	expAt := exp.Unix()
	iat := time.Now().Unix()
	stdClaim := jwt.StandardClaims{
		ExpiresAt: expAt,
		IssuedAt:  iat,
	}
	claim.StandardClaims = stdClaim
	claim.Tokentype = accessToken
	signedToken, err := j.signedToken(&claim)
	if err != nil {
		return nil, err
	}
	jwtToken := &JWTToken{
		Token:    signedToken,
		Claim:    claim,
		ExpireAt: exp,
		Scheme:   bearerScheme,
	}
	return jwtToken, nil
}

func (j *JWTConfig) CreateRefreshToken(claim Claim) (*JWTToken, error) {
	exp := time.Now().Add(j.TokenRefreshExpiration * time.Minute)
	expAt := exp.Unix()
	iat := time.Now().Unix()
	stdClaim := jwt.StandardClaims{
		ExpiresAt: expAt,
		IssuedAt:  iat,
	}
	claim.StandardClaims = stdClaim
	claim.Tokentype = refreshToken
	signedToken, err := j.signedToken(&claim)
	if err != nil {
		return nil, err
	}
	jwtToken := &JWTToken{
		Token:    signedToken,
		Claim:    claim,
		ExpireAt: exp,
		Scheme:   bearerScheme,
	}
	return jwtToken, nil
}

func (j *JWTConfig) signedToken(claim *Claim) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return signedToken, nil
}
