package jwtutil

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SecretKey []byte

// VerifyToken verifies a JWT token and extracts user id from it.
// It takes a JWT token in a string format and returns user id in an int64 format.
// If the token is invalid or expired, it returns an error.
func VerifyToken(ctx context.Context, tokenString string) (int64, error) {

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Check that the signature algorithm matches the expected one.
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return SecretKey, nil
	})
	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	// Check if the token is expired
	if expVal, exists := claims["exp"]; exists {
		switch exp := expVal.(type) {
		case float64:
			if int64(exp) < time.Now().Unix() {
				return 0, fmt.Errorf("token is expired")
			}
		case string:
			parsedExp, err := strconv.ParseInt(exp, 10, 64)
			if err != nil {
				return 0, fmt.Errorf("failed to parse exp from string: %w", err)
			}
			if parsedExp < time.Now().Unix() {
				return 0, fmt.Errorf("token is expired")
			}
		default:
			return 0, fmt.Errorf("unexpected type for exp: %T", expVal)
		}
	} else {
		return 0, fmt.Errorf("exp claim not found in token")
	}

	// For debugging, output the contents of the token to the log
	// TODO: Remove this log
	log.Printf("JWT claims: %+v", claims)

	// Extracting user id from claim "uid"
	uidValue, exists := claims["uid"]
	if !exists {
		return 0, fmt.Errorf("uid not found in token")
	}

	// There are two possible cases float64 or a string.
	switch v := uidValue.(type) {
	case float64:
		return int64(v), nil
	case string:
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse uid from string: %w", err)
		}
		return id, nil
	default:
		return 0, fmt.Errorf("unexpected type for uid: %T", uidValue)
	}
}
