// Copyright 2022 Qian Qiao
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Manager is responsible for creating and validating JWT tokens.
//
// Given that validating JWT comes with a hefty cost, internally, the manager
// caches already validated tokens, so if the same token is validated
// repeatedly, cached results will be returned.
type Manager struct {
	signingKey interface{}
	parseKey   interface{}

	signingMethod jwt.SigningMethod

	validTokens sync.Map
}

// NewManager creates a new JWT client that signs and validates JWT tokens
// using the PS512 algorithm.
//
// Deprecated: Instead of using this method, users of the library should use
// NewPS512Manager instead. The underlying source code has already been
// converted to use the new function, and all users should also do so.
//
// This method will be removed in the 2.0 version stream when we implement
// generics.
func NewManager(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) Manager {
	return NewPS512Manager(publicKey, privateKey)
}

// NewPS512Manager creates a new JWT client that signs and validates JWT tokens
// using the PS512 algorithm.
func NewPS512Manager(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) Manager {
	return Manager{
		parseKey:   publicKey,
		signingKey: privateKey,

		signingMethod: jwt.GetSigningMethod("PS512"),
	}
}

// Alg returns the signing algorithm supported by the current manager instance.
func (m Manager) Alg() string {
	return m.signingMethod.Alg()
}

// CreateCustom creates a JWT token with custom dat claim.
//
// Deprecated: Instead of using this method, users of the library should create
// the Claims object separately and use the SignCustom function instead. The
// underlying source code has already been converted to use the new function,
// and all users should also do so.
//
// This method will be removed in the 2.0 version stream when we implement
// generics.
func (m Manager) CreateCustom(dat interface{},
	expiresAt *time.Time) (<-chan string, <-chan error) {
	_expiresAt := time.Now().Add(365 * 24 * time.Hour)
	if expiresAt != nil {
		_expiresAt = time.Unix(expiresAt.Unix(), 0)
	}

	claims := NewClaims().WithDat(dat).WithExpiry(_expiresAt)

	return m.SignCustom(claims)
}

// ParseCustom parses a JWT token with the claims and returns the claims of
// the token.
//
// TODO: make this generic in 2.0
func (m Manager) ParseCustom(token string) (<-chan *Claims, <-chan error) {
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
						return nil, fmt.Errorf("Unexpected algorithm: %s",
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
				errCh <- errors.New("Unable to parse claims")
				return
			}

			// After parsing the token, we save it to the valid tokens cache
			m.validTokens.Store(token, claims)

			resultCh <- claims
		} else {
			cl := claims.(*Claims)

			if time.Unix(cl.ExpiresAt, 0).Before(time.Now()) {
				errCh <- errors.New("Token Expired")
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
func (m Manager) SignCustom(claims *Claims) (<-chan string, <-chan error) {
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

// ValidateCustom validates a JWT token with a custom dat claim.
//
// Deprecated: Instead of using this method, users of the library should use
// ParseCustom function instead, and subsequently work directly with the Claims
// object that the function returns. The underlying source code has already
// been converted to use the new function, and all users should also do so.
//
// This method will be removed in the 2.0 version stream when we implement
// generics.
func (m Manager) ValidateCustom(token string) (<-chan interface{}, <-chan error) {
	resultCh := make(chan interface{})
	errCh := make(chan error)

	go func() {
		defer close(resultCh)
		defer close(errCh)

		claimsCh, eCh := m.ParseCustom(token)
		select {
		case err := <-eCh:
			errCh <- err
		case claims := <-claimsCh:
			resultCh <- claims.Dat
		}
	}()

	return resultCh, errCh
}
