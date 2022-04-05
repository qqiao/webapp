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
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	f "github.com/qqiao/webapp/firebase/firestore"
)

// Errors
var (
	ErrTokenInvalid = errors.New("token invalid")
)

// Token represents a rememberme token stored in the firestore.
type Token struct {
	Username   string
	Identifier string
	Revoked    bool
	UserAgent  string
	Created    int64
	LastUsed   int64
}

// TokenManager manages datastore operations regarding rememberme tokens.
type TokenManager struct {
	client         *firestore.Client
	collectionName string
}

// NewTokenManager creates a token manager with the given firestore client
// and collection name to store the rememberme tokens in.
func NewTokenManager(client *firestore.Client, collectionName string) TokenManager {
	return TokenManager{
		client:         client,
		collectionName: collectionName,
	}
}

// PurgeTokens removes tokens belonging to a given user last used before or
// equal to the cutoff time.
//
// This function DELETES all matching tokens, regardless of whether the token
// has been revoked.
func (m TokenManager) PurgeTokens(ctx context.Context, username string,
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

			if token, err := iter.Next(); err != nil {
				if err == iterator.Done {
					return nil
				}
				return err
			} else {
				if err = t.Delete(token.Ref); err != nil {
					return err
				}
				return nil
			}
		}); err != nil {
			errCh <- err
		}
	}()

	return errCh
}

// RevokeToken revokes a given token by marking the Revoked field to true.
//
// Although both revoking a token and removing a token will make the
// ValidateToken call fail, RevokeToken leaves the token stored in the data
// store.
func (m TokenManager) RevokeToken(ctx context.Context,
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

			if doc, err := iter.Next(); err != nil {
				return err
			} else {
				m := doc.Data()
				m["Revoked"] = true
				if err = t.Set(doc.Ref, m); err != nil {
					return err
				}
				return nil
			}
		}); err != nil {
			errCh <- err
		}
	}()

	return errCh
}

// SaveToken saves the given rememberme token to the underlying datastore.
func (m TokenManager) SaveToken(ctx context.Context, token Token) <-chan error {
	errCh := make(chan error)

	go func() {
		defer close(errCh)
		now := time.Now().Unix()
		token.Created = now
		token.LastUsed = now
		token.Revoked = false
		doc := m.client.Collection(m.collectionName).NewDoc()
		if _, err := doc.Set(ctx, token); err != nil {
			errCh <- err
		}
	}()

	return errCh
}

// ValidateToken validates if the given token is stored in the datastore.
//
// If the token is valid, that is the token existed, and it has not been
// revoked, its LastUsed will be updated to the current time to record the
// fact that the token has recently been used, otherwise a ErrTokenInvalid will
// be returned.
func (m TokenManager) ValidateToken(ctx context.Context,
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
					{
						Path:     "Revoked",
						Operator: "==",
						Value:    false,
					},
				},
			})

			iter := q.Documents(ctx)
			defer iter.Stop()

			if doc, err := iter.Next(); err != nil {
				if err == iterator.Done {
					return ErrTokenInvalid
				}
				return err
			} else {
				m := doc.Data()
				m["LastUsed"] = time.Now().Unix()
				if err = t.Set(doc.Ref, m); err != nil {
					return err
				}
				return nil
			}
		}); err != nil {
			errCh <- err
		}
	}()

	return errCh
}
