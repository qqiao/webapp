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

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	f "github.com/qqiao/webapp/firebase/firestore"
)

// FirebaseTokenManager manages datastore operations regarding rememberme tokens.
type FirebaseTokenManager struct {
	client         *firestore.Client
	collectionName string
}

// NewTokenManager creates a token manager with the given firestore client
// and collection name to store the rememberme tokens in.
//
// Deprecated: Please use NewFirebaseTokenManager instead.
func NewTokenManager(client *firestore.Client, collectionName string) FirebaseTokenManager {
	return NewFirebaseTokenManager(client, collectionName)
}

// NewFirebaseTokenManager creates a token manager with the given firestore
// client and collection name to store the rememberme tokens in.
func NewFirebaseTokenManager(client *firestore.Client,
	collectionName string) FirebaseTokenManager {
	return FirebaseTokenManager{
		client:         client,
		collectionName: collectionName,
	}
}

// Delete deletes the token permanently from the underlying datastore.
//
// Once deleted, a token cannot be recovered
func (m FirebaseTokenManager) Delete(ctx context.Context,
	token Token) <-chan error {
	errCh := make(chan error)

	go func() {
		defer close(errCh)

		if err := m.client.RunTransaction(ctx, func(ctx context.Context,
			t *firestore.Transaction) error {
			q := f.ApplyQuery(m.client.Collection(m.collectionName), f.Query{
				Filters: []f.Filter{
					{
						Path:     "Username",
						Operator: "==",
						Value:    token.Username,
					},
					{
						Path:     "Identifier",
						Operator: "==",
						Value:    token.Identifier,
					},
				},
			})

			iter := q.Documents(ctx)
			defer iter.Stop()

			for {
				ds, err := iter.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					return err
				}

				if err = t.Delete(ds.Ref); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			errCh <- err
		}
	}()

	return errCh
}

// Purge removes tokens belonging to a given user last used before or equal to
// the cutoff time.
//
// This function DELETES all matching tokens, regardless of whether the token
// has been revoked.
func (m FirebaseTokenManager) Purge(ctx context.Context, username string,
	cutoff time.Time) <-chan error {
	errCh := make(chan error)

	go func() {
		defer close(errCh)

		if err := m.client.RunTransaction(ctx, func(ctx context.Context,
			t *firestore.Transaction) error {
			q := f.ApplyQuery(m.client.Collection(m.collectionName), f.Query{
				Filters: []f.Filter{
					{
						Path:     "Username",
						Operator: "==",
						Value:    username,
					},
					{
						Path:     "LastUsed",
						Operator: "<=",
						Value:    cutoff.Unix(),
					},
				},
			})

			iter := q.Documents(ctx)
			defer iter.Stop()

			for {
				ds, err := iter.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					return err
				}

				if err = t.Delete(ds.Ref); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			errCh <- err
		}
	}()

	return errCh
}

// PurgeTokens removes tokens belonging to a given user last used before or
// equal to the cutoff time.
//
// This function DELETES all matching tokens, regardless of whether the token
// has been revoked.
//
// Deprecated: please use method Purge instead
func (m FirebaseTokenManager) PurgeTokens(ctx context.Context, username string,
	cutoff time.Time) <-chan error {
	return m.Purge(ctx, username, cutoff)
}

// Revoke revokes a given token by marking the Revoked field to true.
//
// Although both revoking a token and removing a token will make the
// ValidateToken call fail, RevokeToken leaves the token stored in the data
// store.
func (m FirebaseTokenManager) Revoke(ctx context.Context,
	token Token) (<-chan *Token, <-chan error) {
	tokenCh := make(chan *Token)
	errCh := make(chan error)

	go func() {
		defer close(tokenCh)
		defer close(errCh)

		var tok Token

		if err := m.client.RunTransaction(ctx, func(ctx context.Context,
			t *firestore.Transaction) error {
			q := f.ApplyQuery(m.client.Collection(m.collectionName), f.Query{
				Filters: []f.Filter{
					{
						Path:     "Username",
						Operator: "==",
						Value:    token.Username,
					},
					{
						Path:     "Identifier",
						Operator: "==",
						Value:    token.Identifier,
					},
				},
			})

			iter := q.Documents(ctx)
			defer iter.Stop()

			ds, err := iter.Next()
			if err != nil {
				return err
			}

			if err = ds.DataTo(&tok); err != nil {
				return err
			}

			tok.Revoked = true
			if err = t.Set(ds.Ref, tok); err != nil {
				return err
			}
			return nil

		}); err != nil {
			errCh <- err
		}
		tokenCh <- &tok
	}()

	return tokenCh, errCh
}

// RevokeToken revokes a given token by marking the Revoked field to true.
//
// Although both revoking a token and removing a token will make the
// ValidateToken call fail, RevokeToken leaves the token stored in the data
// store.
//
// Deprecated: please use Revoke instead
func (m FirebaseTokenManager) RevokeToken(ctx context.Context,
	token Token) <-chan error {
	t, e := m.Revoke(ctx, token)
	errCh := make(chan error)

	go func() {
		defer close(errCh)

		select {
		case err := <-e:
			errCh <- err

		case <-t:
		}
	}()
	return errCh
}

// Save saves the token to the underlying datastore
//
// This function will return ErrTokenDuplicate if the given Username
// Identifier combination already exists in the datastore
func (m FirebaseTokenManager) Save(ctx context.Context,
	token Token) (<-chan *Token, <-chan error) {
	tokenCh := make(chan *Token)
	errCh := make(chan error)

	go func() {
		defer close(tokenCh)
		defer close(errCh)

		now := time.Now().Unix()
		token.Created = now
		token.LastUsed = now
		token.Revoked = false

		if err := m.client.RunTransaction(ctx, func(ctx context.Context,
			t *firestore.Transaction) error {
			q := f.ApplyQuery(m.client.Collection(m.collectionName), f.Query{
				Filters: []f.Filter{
					{
						Path:     "Username",
						Operator: "==",
						Value:    token.Username,
					},
					{
						Path:     "Identifier",
						Operator: "==",
						Value:    token.Identifier,
					},
				},
			})

			iter := q.Documents(ctx)
			defer iter.Stop()

			_, err := iter.Next()
			if err != iterator.Done {
				return ErrTokenDuplicate
			}

			doc := m.client.Collection(m.collectionName).NewDoc()
			if err = t.Set(doc, token); err != nil {
				return err
			}
			return nil

		}); err != nil {
			errCh <- err
		}

		tokenCh <- &token
	}()

	return tokenCh, errCh
}

// SaveToken saves the given rememberme token to the underlying datastore.
//
// Deprecated: please use the Save method instead
func (m FirebaseTokenManager) SaveToken(ctx context.Context,
	token Token) <-chan error {
	t, e := m.Save(ctx, token)
	errCh := make(chan error)

	go func() {
		defer close(errCh)

		select {
		case err := <-e:
			errCh <- err

		case <-t:
		}
	}()
	return errCh
}

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
// If the token is valid, its LastUsed will be updated to the current time to
// record the fact that the token has recently been used.
func (m FirebaseTokenManager) Validate(ctx context.Context,
	token Token) (<-chan *Token, <-chan error) {
	tokenCh := make(chan *Token)
	errCh := make(chan error)

	go func() {
		defer close(tokenCh)
		defer close(errCh)

		var tok Token

		if err := m.client.RunTransaction(ctx, func(ctx context.Context,
			t *firestore.Transaction) error {
			q := f.ApplyQuery(m.client.Collection(m.collectionName), f.Query{
				Filters: []f.Filter{
					{
						Path:     "Username",
						Operator: "==",
						Value:    token.Username,
					},
					{
						Path:     "Identifier",
						Operator: "==",
						Value:    token.Identifier,
					},
					{
						Path:     "Revoked",
						Operator: "==",
						Value:    false,
					},
				},
			})

			iter := q.Documents(ctx)
			defer iter.Stop()

			ds, err := iter.Next()
			if err == iterator.Done {
				return ErrTokenInvalid
			}
			if err != nil {
				return err
			}

			if err = ds.DataTo(&tok); err != nil {
				return err
			}
			tok.LastUsed = time.Now().Unix()
			if err = t.Set(ds.Ref, tok); err != nil {
				return err
			}
			return nil

		}); err != nil {
			errCh <- err
			return
		}

		tokenCh <- &tok
	}()

	return tokenCh, errCh
}

// ValidateToken validates if the given token is stored in the datastore.
//
// If the token is valid, that is the token existed, and it has not been
// revoked, its LastUsed will be updated to the current time to record the
// fact that the token has recently been used, otherwise a ErrTokenInvalid will
// be returned.
//
// Deprecated: please use the Validate method instead
func (m FirebaseTokenManager) ValidateToken(ctx context.Context,
	token Token) <-chan error {
	errCh := make(chan error)

	go func() {
		defer close(errCh)
		t, e := m.Validate(ctx, token)

		select {
		case err := <-e:
			errCh <- err
		case <-t:
		}
	}()

	return errCh
}
