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
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims represents a custom claim where the dat section is used for custom
// data.
//
// TODO: make this generic in 2.0
type Claims struct {
	Dat interface{} `json:"dat,omitempty"`
	*jwt.StandardClaims
}

// NewClaims creates a new instance of the custom JWT claims.
//
// TODO: make this generic in 2.0
func NewClaims() *Claims {
	return &Claims{
		StandardClaims: &jwt.StandardClaims{},
	}
}

// WithDat adds a dat claim to the JWT token.
//
// TODO: make this generic in 2.0
func (c *Claims) WithDat(dat interface{}) *Claims {
	c.Dat = dat
	return c
}

// WithExpiry updates the expiry of the JWT token to the time specified.
//
// TODO: make this generic in 2.0
func (c *Claims) WithExpiry(expiry time.Time) *Claims {
	c.StandardClaims.ExpiresAt = expiry.Unix()
	return c
}
