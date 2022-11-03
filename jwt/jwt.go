package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claim struct {
	jwt.StandardClaims
	TokenId   string `json:"jti"`
	TokenType string `json:"token_type"`
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
	ExpectAccessToken(token string) (*Claim, error)
	ExpectRefreshToken(token string) (*Claim, error)
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
	claim.TokenType = accessToken
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
	claim.TokenType = refreshToken
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

type JWTError string

func (e JWTError) Error() string {
	return string(e)
}

const JWTExpirationError = JWTError("token is expired")

func (j *JWTConfig) ExpectAccessToken(token string) (*Claim, error) {
	c := &Claim{}
	_, err := jwt.ParseWithClaims(token, c, func(t *jwt.Token) (interface{}, error) {
		if _, isValid := t.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		validationErr, ok := err.(*jwt.ValidationError)
		if ok {
			if validationErr.Errors == jwt.ValidationErrorExpired {
				return nil, JWTExpirationError
			}
		}
		return nil, fmt.Errorf("failed ParseWithClaims: %w", err)
	}
	if c.UserId == 0 {
		return nil, fmt.Errorf("invalid user_id claim")
	}
	if c.TokenType != accessToken {
		return nil, fmt.Errorf("invalid token_type claim")
	}

	return c, nil
}

func (j *JWTConfig) ExpectRefreshToken(token string) (*Claim, error) {
	c := &Claim{}
	_, err := jwt.ParseWithClaims(token, c, func(t *jwt.Token) (interface{}, error) {
		if _, isValid := t.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("unexpected signin method: %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		validationErr, ok := err.(*jwt.ValidationError)
		if ok {
			if validationErr.Errors == jwt.ValidationErrorExpired {
				return nil, JWTExpirationError
			}
		}
		return nil, fmt.Errorf("failed ParseWithClaims: %w", err)
	}
	if c.TokenId == "" {
		return nil, fmt.Errorf("invalid jti claim")
	}
	if c.UserId == 0 {
		return nil, fmt.Errorf("invalid user_id claim")
	}
	if c.TokenType != refreshToken {
		return nil, fmt.Errorf("invalid token_type claim")
	}
	return c, nil
}
