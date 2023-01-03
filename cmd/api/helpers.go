package main

import (
	"github.com/pascaldekloe/jwt"
	"strconv"
	"time"
)

func generate_jwt_token(userID int, baseURL string, secretKey string) (map[string]string, error) {
	var claims jwt.Claims
	claims.Subject = strconv.Itoa(userID)

	expiry := time.Now().Add(24 * time.Hour)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(expiry)

	claims.Issuer = baseURL
	claims.Audiences = []string{baseURL}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(secretKey))
	if err != nil {
		return nil, err
	}

	data := map[string]string{
		"AuthenticationToken":       string(jwtBytes),
		"AuthenticationTokenExpiry": expiry.Format(time.RFC3339),
	}
	return data, nil
}
