package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims represents a custom claim where the dat section is used for custom
// data.
//
// TODO make this generic
type Claims struct {
	Dat interface{} `json:"dat,omitempty"`
	*jwt.StandardClaims
}

// Manager is responsible for creating and validating JWT tokens.
//
// Given that validating JWT comes with a cost, internally, the manager
// caches already validated tokens, so if the same token is validated again
// multiple times, cached results will be returned.
type Manager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey

	validTokens sync.Map
}

// NewManager creates a new JWT client that signs and validates JWT tokens.
//
// TODO make this generic
func NewManager(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) Manager {
	return Manager{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

// CreateCustom creates a JWT token with custom dat claim.
//
// TODO: Make this function generic
func (m Manager) CreateCustom(dat interface{},
	expiresAt *time.Time) (<-chan string, <-chan error) {
	tokenCh := make(chan string)
	errCh := make(chan error)

	go func() {
		defer close(errCh)
		defer close(tokenCh)

		_expiresAt := time.Now().Add(365 * 24 * time.Hour)
		if expiresAt == nil {
			expiresAt = &_expiresAt
		}
		claim := Claims{
			dat,
			&jwt.StandardClaims{
				ExpiresAt: expiresAt.Unix(),
			},
		}

		if token, err := jwt.NewWithClaims(jwt.SigningMethodPS512,
			claim).SignedString(m.privateKey); err != nil {
			errCh <- err
			return
		} else {
			tokenCh <- token
		}
	}()

	return tokenCh, errCh
}

// ValidateCustom validates a JWT token with a custom dat claim.
//
// TODO make this generic
func (m Manager) ValidateCustom(token string) (<-chan interface{}, <-chan error) {
	resultCh := make(chan interface{})
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
					if _, ok := jwtToken.Method.(*jwt.SigningMethodRSAPSS); !ok {
						return nil, fmt.Errorf("Unexpected method: %s",
							jwtToken.Header["alg"])
					}

					return m.publicKey, nil
				})
			if err != nil {
				errCh <- err
				return
			}

			claims, ok := t.Claims.(*Claims)
			if !ok {
				errCh <- errors.New("Unable to parse claims")
				return
			}

			// After parsing the token, we save it to the valid tokens cache
			m.validTokens.Store(token, claims)

			resultCh <- claims.Dat
		} else {
			cl := claims.(*Claims)

			if time.Unix(cl.ExpiresAt, 0).Before(time.Now()) {
				errCh <- errors.New("Token Expired")
				return
			}

			resultCh <- cl.Dat
		}
	}()

	return resultCh, errCh
}
