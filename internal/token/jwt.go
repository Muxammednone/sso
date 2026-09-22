package token

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTGenerator struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTGenerator(
	secret []byte,
	ttl time.Duration,
) *JWTGenerator {
	return &JWTGenerator{
		secret: secret,
		ttl:    ttl,
	}
}

func (g *JWTGenerator) Generate(
	userID int64,
) (string, error) {
	now := time.Now()

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.FormatInt(userID, 10),
		"iat": now.Unix(),
		"exp": now.Add(g.ttl).Unix(),
	})

	return t.SignedString(g.secret)
}
