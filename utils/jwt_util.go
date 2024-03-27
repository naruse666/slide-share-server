package utils

import (
	"fmt"
	"slide-share/model"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

func VerifyAndGetUserClaims(authToken string, secret string) (*model.JWTPayload, error) {
	token, err := jwt.Parse(strings.TrimPrefix(authToken, "Bearer "), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		println("error", err.Error())
		return nil, err
	}

	var payload model.JWTPayload
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, ok := claims["id"].(string)
		if !ok || id == "" {
			return nil, fmt.Errorf("id is empty or not a string")
		}
		role, ok := claims["role"].(string)
		if !ok || role == "" {
			return nil, fmt.Errorf("role is empty or not a string")
		}
		payload = model.JWTPayload{
			ID:   id,
			Role: role,
		}
		return &payload, nil
	}

	return nil, fmt.Errorf("invalid token")
}
