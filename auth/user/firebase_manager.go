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
	"github.com/qqiao/webapp/datastore"
	f "github.com/qqiao/webapp/firebase/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const defaultBatchSize = 10

// FirebaseManager is responsible for all user related operations
type FirebaseManager struct {
	client         *firestore.Client
	collectionName string
}

// NewFirebaseManager creates a new UserManager with the given firestore client
// and collection name.
func NewFirebaseManager(client *firestore.Client,
	collectionName string) FirebaseManager {
	return FirebaseManager{
		client:         client,
		collectionName: collectionName,
	}
}

// Add adds a user to the database of users.
//
// Please note that Add will return ErrUserDuplicate if the user already exists
// in the datastore.
func (m FirebaseManager) Add(ctx context.Context,
	usr User) (<-chan *User, <-chan error) {
	userCh := make(chan *User)
	errCh := make(chan error)

	go func() {
		defer close(userCh)
		defer close(errCh)

		if err := m.client.RunTransaction(ctx, func(ctx context.Context,
			t *firestore.Transaction) error {
			ref := m.client.Collection(m.collectionName).Doc(usr.Username)
			_, err := t.Get(ref)
			if err == nil {
				return ErrUserDuplicate
			}

			if err != nil && status.Code(err) != codes.NotFound {
				return err
			}

			return t.Set(ref, usr)
		}); err != nil {
			errCh <- err
			return
		}
		userCh <- &usr
	}()

	return userCh, errCh
}

// ConfirmExists looks for the user with the given username and password.
//
// This function will return ErrUserNotFound in the case where a user
// matching the given user cannot be found.
func (m FirebaseManager) ConfirmExists(ctx context.Context, user User) <-chan error {
	errCh := make(chan error)

	go func() {
		defer close(errCh)

		q := f.ApplyQuery(m.client.Collection(m.collectionName),
			datastore.Query{
				Filters: []datastore.Filter{
					{
						Path:     "Username",
						Operator: "==",
						Value:    user.Username,
					},
					{
						Path:     "Password",
						Operator: "==",
						Value:    user.Password,
					},
				},
			})

		iter := q.Documents(ctx)
		defer iter.Stop()

		if _, err := iter.Next(); err != nil {
			// If we immediately get the done error, it means that we didn't find
			// the user
			if err == iterator.Done {
				errCh <- ErrUserNotFound
			} else {
				errCh <- err
			}
		}
	}()

	return errCh
}

// Find finds the user based on the given query criteron
func (m FirebaseManager) Find(ctx context.Context,
	query datastore.Query) (<-chan (<-chan *User), <-chan error) {

	resultsCh := make(chan (<-chan *User))
	errCh := make(chan error)

	go func() {
		defer close(resultsCh)
		defer close(errCh)

		// We check the batch size here so that we can use a buffered channel for
		// better performance
		batchSize := defaultBatchSize
		if query.Limit != 0 {
			batchSize = query.Limit
		}
		usersCh := make(chan *User, batchSize)
		defer close(usersCh)

		q := f.ApplyQuery(m.client.Collection(m.collectionName), query)

		iter := q.Documents(ctx)
		defer iter.Stop()

		for {
			ref, err := iter.Next()
			if err == iterator.Done {
				resultsCh <- usersCh
				return
			}

			if err != nil {
				errCh <- err
				return
			}

			var user User
			if err = ref.DataTo(&user); err != nil {
				errCh <- err
				return
			}

			usersCh <- &user
		}
	}()

	return resultsCh, errCh
}

// Update updates the given user record.
//
// Update will return ErrUserNotFound if the user cannot be found in the
// underlying datastore
func (m FirebaseManager) Update(ctx context.Context,
	user User) (<-chan *User, <-chan error) {
	userCh := make(chan *User)
	errCh := make(chan error)

	go func() {
		defer close(userCh)
		defer close(errCh)

		if err := m.client.RunTransaction(ctx, func(ctx context.Context,
			t *firestore.Transaction) error {
			ref := m.client.Collection(m.collectionName).Doc(user.Username)
			_, err := t.Get(ref)
			if err != nil && status.Code(err) != codes.NotFound {
				return ErrUserNotFound
			}

			return t.Set(ref, user)
		}); err != nil {
			errCh <- err
			return
		}
		userCh <- &user
	}()

	return userCh, errCh
}
