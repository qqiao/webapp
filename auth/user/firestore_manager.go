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

package user

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/qqiao/webapp/v2/datastore"
	f "github.com/qqiao/webapp/v2/datastore/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FirestoreManager is an UserManager implementation that uses firebase
// firestore as the underlying user information storage engine.
type FirestoreManager struct {
	client         *firestore.Client
	collectionName string
}

// NewFirestoreManager creates a new FirestoreManager with the given firestore
// client and collection name.
func NewFirestoreManager(client *firestore.Client,
	collectionName string) *FirestoreManager {
	return &FirestoreManager{
		client:         client,
		collectionName: collectionName,
	}
}

// Add adds a user to the database of users.
//
// Please note that a user is considered a duplicate if any of the following
// already exist on a different user: Email, PhoneNumber, and Username. The Add
// method will return ErrUserDuplicate in this case.
func (m *FirestoreManager) Add(ctx context.Context, usr *User) (<-chan *User,
	<-chan error) {
	userCh := make(chan *User)
	errCh := make(chan error)

	go func() {
		defer close(userCh)
		defer close(errCh)

		if err := m.client.RunTransaction(ctx, func(ctx context.Context,
			t *firestore.Transaction) error {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			col := m.client.Collection(m.collectionName)

			queries := make([]datastore.Query, 0)
			if usr.Username != "" {
				queries = append(queries, datastore.Query{
					Filters: []datastore.Filter{{
						Path:     "Username",
						Operator: "==",
						Value:    usr.Username,
					}},
				})
			}
			if usr.Email != "" {
				queries = append(queries, datastore.Query{
					Filters: []datastore.Filter{{
						Path:     "Email",
						Operator: "==",
						Value:    usr.Email,
					}},
				})
			}
			if usr.PhoneNumber != "" {
				queries = append(queries, datastore.Query{
					Filters: []datastore.Filter{{
						Path:     "PhoneNumber",
						Operator: "==",
						Value:    usr.PhoneNumber,
					}},
				})
			}

			found, errs := f.Or[*User](ctx, 5, 5, t, col, queries...)
			for done := false; !done; {
				select {
				case u, ok := <-found:
					if ok && u != nil {
						cancel()
						return ErrUserDuplicate
					}
					done = true
				case er, ok := <-errs:
					if ok && er != nil {
						cancel()
						return er
					}
				}
			}

			ref := col.NewDoc()
			usr.UID = ref.ID
			return t.Set(ref, usr)
		}); err != nil {
			errCh <- err
			return
		}
		userCh <- usr
	}()

	return userCh, errCh
}

// Find finds the user based on the given query criterion.
//
// If multiple queries are sent, the queries are combined with OR
// condition. Please refer to https://pkg.go.dev/github.com/qqiao/webapp/v2/datastore/firestore#Or
// for limitations of OR queries.
func (m *FirestoreManager) Find(ctx context.Context,
	queries ...datastore.Query) (<-chan *User, <-chan error) {
	out := make(chan *User)
	errs := make(chan error)

	go func() {
		defer close(out)
		defer close(errs)

		if err := m.client.RunTransaction(ctx, func(ctx context.Context,
			t *firestore.Transaction) error {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			col := m.client.Collection(m.collectionName)

			found, errCh := f.Or[*User](ctx, 5, 5, t, col, queries...)
			for done := false; !done; {
				select {
				case u, ok := <-found:
					if ok && u != nil {
						out <- u
					} else {
						done = true
					}
				case err, ok := <-errCh:
					if ok && err != nil {
						errs <- err
					}
				}
			}
			return nil
		}); err != nil {
			errs <- err
			return
		}
	}()

	return out, errs
}

// Update updates the given user record.
//
// Update will return ErrUserNotFound if the user cannot be found in the
// underlying datastore
func (m *FirestoreManager) Update(ctx context.Context,
	usr *User) (<-chan *User, <-chan error) {
	userCh := make(chan *User)
	errCh := make(chan error)

	go func() {
		defer close(userCh)
		defer close(errCh)

		if err := m.client.RunTransaction(ctx, func(ctx context.Context,
			t *firestore.Transaction) error {
			ref := m.client.Collection(m.collectionName).Doc(usr.UID)
			_, err := t.Get(ref)
			if err != nil {
				if status.Code(err) == codes.NotFound {
					return ErrUserNotFound
				}
				return err
			}

			return t.Set(ref, usr)
		}); err != nil {
			errCh <- err
			return
		}
		userCh <- usr
	}()

	return userCh, errCh
}
