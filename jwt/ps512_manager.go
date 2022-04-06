package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// PS512Manager is responsible for creating and validating JWT tokens using
// PS512 algorithm.
//
// Given that validating JWT comes with a hefty cost, internally, the manager
// caches already validated tokens, so if the same token is validated
// repeatedly, cached results will be returned.
type PS512Manager struct {
	signingKey interface{}
	parseKey   interface{}

	signingMethod jwt.SigningMethod

	validTokens sync.Map
}

// NewPS512Manager creates a new JWT client that signs and validates JWT tokens
// using the PS512 algorithm.
func NewPS512Manager(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) PS512Manager {
	return PS512Manager{
		parseKey:   publicKey,
		signingKey: privateKey,

		signingMethod: jwt.GetSigningMethod("PS512"),
	}
}

// Alg returns the signing algorithm supported by the current manager instance.
func (m *PS512Manager) Alg() string {
	return m.signingMethod.Alg()
}

// ParseCustom parses a JWT token with the claims and returns the claims of
// the token.
//
// TODO: make this generic in 2.0
func (m *PS512Manager) ParseCustom(token string) (<-chan *Claims,
	<-chan error) {
	resultCh := make(chan *Claims)
	errCh := make(chan error)

	go func() {
		defer close(resultCh)
		defer close(errCh)

		// First let check if we have it in the valid tokens cache
		claims, has := m.validTokens.Load(token)

		// If we don't, we parse the token
		if !has {
			t, err := jwt.ParseWithClaims(token, &Claims{},
				func(jwtToken *jwt.Token) (interface{}, error) {
					if jwtToken.Method.Alg() != m.Alg() {
						return nil, fmt.Errorf("unexpected algorithm: %s",
							jwtToken.Header["alg"])
					}

					return m.parseKey, nil
				})
			if err != nil {
				errCh <- err
				return
			}

			claims, ok := t.Claims.(*Claims)
			if !ok {
				errCh <- errors.New("unable to parse claims")
				return
			}

			// After parsing the token, we save it to the valid tokens cache
			m.validTokens.Store(token, claims)

			resultCh <- claims
		} else {
			cl := claims.(*Claims)

			if time.Unix(cl.ExpiresAt, 0).Before(time.Now()) {
				errCh <- errors.New("token expired")
				return
			}

			resultCh <- cl
		}
	}()

	return resultCh, errCh
}

// SignCustom signs the JWT token with the given claims.
//
// TODO: make this generic in 2.0
func (m *PS512Manager) SignCustom(claims *Claims) (<-chan string,
	<-chan error) {
	tokenCh := make(chan string)
	errCh := make(chan error)

	go func() {
		defer close(errCh)
		defer close(tokenCh)

		if token, err := jwt.NewWithClaims(m.signingMethod, claims).
			SignedString(m.signingKey); err != nil {
			errCh <- err
			return
		} else {
			tokenCh <- token
		}
	}()

	return tokenCh, errCh
}
