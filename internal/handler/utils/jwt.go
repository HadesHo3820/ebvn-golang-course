package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrEmptyUID     = errors.New("empty uid")
)

// GetJWTClaimsFromContext returns jwt.MapClaims from request's token
func GetJWTClaimsFromRequest(c *gin.Context) (jwt.MapClaims, error) {
	tokenInfo, _ := c.Get("claims")
	claims, valid := tokenInfo.(jwt.MapClaims)
	if !valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GetUIDFromRequest returns uid from request's token
func GetUIDFromRequest(c *gin.Context) (string, error) {
	claims, err := GetJWTClaimsFromRequest(c)
	if err != nil {
		return "", err
	}

	uid, ok := claims["sub"].(string)
	if !ok || uid == "" {
		return "", ErrEmptyUID
	}

	return uid, nil
}
