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

// Manager is responsible for all the JWT token related operations.
type Manager interface {
	// Alg returns the signing algorithm supported by the current manager
	// instance.
	Alg() string

	// ParseCustom parses a JWT token with the claims and returns the claims of
	// the token.
	ParseCustom(token string) (<-chan *Claims, <-chan error)

	// SignCustom signs the JWT token with the given claims.
	SignCustom(claims *Claims) (<-chan string, <-chan error)
}
