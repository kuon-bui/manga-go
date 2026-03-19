package jwtprovider

import (
	"context"
	"errors"
	"fmt"
	"manga-go/internal/pkg/config"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/redis"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserPayload struct {
	UserID       uuid.UUID `json:"userId"`
	FullName     string    `json:"fullName"`
	Email        string    `json:"email"`
	IsRefresh    bool      `json:"isRefresh"`
	RefreshCount int       `json:"refreshCount"`
}

type JwtProvider struct {
	config           *config.Config
	logger           *logger.Logger
	rds              *redis.Redis
	jwtSecret        []byte
	jwtRefreshSecret []byte
}

type CustomClaims struct {
	jwt.RegisteredClaims
	UserPayload
}

const blacklistTokenPrefix = "blacklist_token:"

func NewJwtProvider(config *config.Config, logger *logger.Logger, rds *redis.Redis) *JwtProvider {
	return &JwtProvider{
		config:           config,
		logger:           logger,
		rds:              rds,
		jwtSecret:        []byte(config.Jwt.Secret),
		jwtRefreshSecret: []byte(config.Jwt.RefreshSecret),
	}
}

func (p *JwtProvider) generate(userPayload UserPayload, oldTokenId string) (*Token, *Token, error) {
	now := time.Now()

	id := oldTokenId
	if id == "" {
		id = uuid.NewString()
	}

	accessTokenExpire := now.Add(time.Duration(p.config.Jwt.ExpiresAt) * time.Second).Add(7 * time.Hour)
	refreshTokenExpire := now.Add(time.Duration(p.config.Jwt.RefreshExpire) * time.Second).Add(7 * time.Hour)

	registeredClaims := jwt.RegisteredClaims{
		ID:        id,
		ExpiresAt: jwt.NewNumericDate(accessTokenExpire),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now.Add(-time.Nanosecond)),
	}
	claims := CustomClaims{
		RegisteredClaims: registeredClaims,
		UserPayload:      userPayload,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(p.jwtSecret)
	if err != nil {
		return nil, nil, err
	}

	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(refreshTokenExpire)

	claims.UserPayload.IsRefresh = true
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshTokenString, err := refreshToken.SignedString(p.jwtRefreshSecret)
	if err != nil {
		return nil, nil, err
	}

	return &Token{TokenString: tokenString, ExpiresAt: accessTokenExpire},
		&Token{TokenString: refreshTokenString, ExpiresAt: refreshTokenExpire},
		err
}

type Token struct {
	TokenString string
	ExpiresAt   time.Time
}

// GenerateToken generates a JWT token and a refresh token for the given user payload.
// The tokens are signed using the secret and expiration time defined in the configuration.
// Returns the generated token, refresh token, and any error encountered.
func (p *JwtProvider) GenerateToken(userPayload UserPayload) (*Token, *Token, error) {
	return p.generate(userPayload, "")
}

func (p *JwtProvider) validate(ctx context.Context, tokenString string, isRefresh bool) (*UserPayload, string, error) {
	tokenInvalid := errors.New("claims not valid")

	keyFunc := func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		if isRefresh {
			return p.jwtRefreshSecret, nil
		}

		return p.jwtSecret, nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, keyFunc)
	if err != nil {
		return nil, "", err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, "", tokenInvalid
	}

	isBlock, err := p.checkBlockTokenId(ctx, claims.ID)
	if err != nil {
		return nil, "", err
	}

	if isBlock {
		return nil, "", tokenInvalid
	}

	if isRefresh != claims.UserPayload.IsRefresh {
		return nil, "", tokenInvalid
	}

	return &claims.UserPayload, claims.ID, nil
}

func (p *JwtProvider) ValidateToken(ctx context.Context, tokenString string) (*UserPayload, string, error) {
	return p.validate(ctx, tokenString, false)
}

func (p *JwtProvider) RenewAccessToken(ctx context.Context, refreshTokenString string) (*Token, *Token, error) {
	userPayload, tokenID, err := p.validate(ctx, refreshTokenString, true)
	if err != nil {
		return nil, nil, err
	}

	userPayload.RefreshCount++
	userPayload.IsRefresh = false

	return p.generate(*userPayload, tokenID)
}

func (p *JwtProvider) checkBlockTokenId(ctx context.Context, tokenID string) (bool, error) {
	isBlocked, err := p.rds.Client().Get(ctx, blacklistTokenPrefix+tokenID).Bool()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}

		return false, err
	}
	return isBlocked, nil
}

func (p *JwtProvider) InvalidateToken(ctx context.Context, tokenId string) error {
	return p.rds.Client().Set(ctx, blacklistTokenPrefix+tokenId, true, time.Duration(p.config.Jwt.RefreshExpire)*time.Second).Err()
}
