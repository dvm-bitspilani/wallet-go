package helpers

import (
	"github.com/pascaldekloe/jwt"
	"strconv"
	"time"
)

func GenerateJwtToken(userID int, baseURL string, secretKey string) (map[string]string, error) {
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

type OrderActionOrderStruct struct {
	ItemId   int `json:"item_id"`
	Quantity int `json:"quantity"`
}

type OrderActionVendorStruct struct {
	VendorId int                      `json:"vendor_id"`
	Order    []OrderActionOrderStruct `json:"order"`
}
