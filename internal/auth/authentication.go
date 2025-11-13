package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	hashedPass, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	return hashedPass, nil
}

func CheckHashedPassword(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return match, err
	}

	return match, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "PortfolioTrackerV2",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})

	signed, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok || t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(tokenSecret), nil
	})
	if err != nil || !token.Valid {
		return uuid.Nil, err
	}
	return uuid.Parse(claims.Subject)
}

func GetBearerToken(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	var token string
	if len(val) == 0 {
		return "", fmt.Errorf("header not found")
	}

	valTrimmed := strings.TrimSpace(val)
	valSplit := strings.SplitN(valTrimmed, " ", 2)

	if valSplit[0] != "Bearer" && len(valSplit[1]) != 0 {
		token = strings.TrimSpace(valSplit[1])
	} else {
		return "", fmt.Errorf("unable to extract token")
	}

	return token, nil
}
