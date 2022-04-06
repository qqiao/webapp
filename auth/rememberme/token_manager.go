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

package rememberme

import (
	"context"
	"time"
)

// TokenManager manages all rememberme token related operations.
type TokenManager interface {
	// Delete deletes the token permanently from the underlying datastore.
	//
	// Once deleted, a token cannot be recovered
	Delete(ctx context.Context, token Token) <-chan error

	// Purge removes tokens belonging to a given user last used before or equal
	// to the cutoff time.
	//
	// This function DELETES all matching tokens, regardless of whether the
	// token has been revoked.
	Purge(ctx context.Context, username string, cutoff time.Time) <-chan error

	// Revoke revokes a given token by marking the Revoked field to true.
	//
	// Although both revoking a token and removing a token will make the
	// ValidateToken call fail, RevokeToken leaves the token stored in the data
	// store.
	Revoke(ctx context.Context, token Token) (<-chan *Token, <-chan error)

	// Save saves the token to the underlying datastore
	//
	// This function will return ErrTokenDuplicate if the given Username
	// Identifier combination already exists in the datastore
	Save(ctx context.Context, token Token) (<-chan *Token, <-chan error)

	// Validate checks if the given token is valid.
	//
	// A token is considered valid if it meets the following conditions:
	//
	//   1. The Username/Identifier combination exists in the datastore
	//   2. The token has not been revoked.
	//
	// This method returns a ErrTokenInvalid if the token cannot be validated.
	// This method also passes through any underlying datastore errors to the
	// caller.
	//
	// If the token is valid, its LastUsed will be updated to the current time
	// to record the fact that the token has recently been used.
	Validate(ctx context.Context, token Token) (<-chan *Token, <-chan error)
}
