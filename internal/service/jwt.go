package service

import (
	"time"

	"github.com/RethikRaj/task_manager_go/internal/errs"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	UserID int `json:"sub"`
	// Role   string `json:"role"` TODO:later
	jwt.RegisteredClaims // Embeds standard fields like exp, iat, iss
}

func GenerateToken(userId int, secret string) (string, error) {
	claims := UserClaims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string, secret string) (*UserClaims, error) {
	claims := &UserClaims{}

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		// Validate whether the algorithm is what we expect
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.ErrUnexpectedSigningMethod
		}

		// Return the secret key used for signing
		return []byte(secret), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)

	if err != nil {
		return nil, err // Could be expired, malformed, or invalid signature
	}

	// Final check: is the token valid?
	if !token.Valid {
		return nil, errs.ErrInvalidToken
	}

	return claims, nil
}
