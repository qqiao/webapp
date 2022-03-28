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
	ErrTokenInvalid = errors.New("Token invalid")
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
func (m TokenManager) PurgeTokens(ctx context.Context, username string, cutoff time.Time) <-chan error {
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
func (m TokenManager) RevokeToken(ctx context.Context, token Token) <-chan error {
	errCh := make(chan error)

	go func() {
		defer close(errCh)

		// TODO implement this
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
// If the token is valid, that is the token existed and it has not been
// revoked, its LastUsed will be updated to the current time to record the
// fact that the token has recently be used, otherwise a ErrTokenInvalid will
// be returned.
func (m TokenManager) ValidateToken(ctx context.Context, token Token) <-chan error {
	errCh := make(chan error)

	go func() {
		defer close(errCh)

		// TODO implement this
	}()

	return errCh
}
