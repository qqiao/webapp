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
	"github.com/qqiao/pipeline/v2"
	"github.com/qqiao/webapp/v2/datastore"
	f "github.com/qqiao/webapp/v2/datastore/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const defaultBatchSize = 10

// FirestoreManager is responsible for all user related operations
type FirestoreManager struct {
	client         *firestore.Client
	collectionName string
}

// NewFirestoreManager creates a new UserManager with the given firestore client
// and collection name.
func NewFirestoreManager(client *firestore.Client,
	collectionName string) *FirestoreManager {
	return &FirestoreManager{
		client:         client,
		collectionName: collectionName,
	}
}

func (m *FirestoreManager) createStreamWorker(t *firestore.Transaction) pipeline.
	StreamWorker[any, any] {
	return func(ctx context.Context, in pipeline.Producer[any]) (
		<-chan any,
		<-chan error) {
		results := make(chan any)
		errCh := make(chan error)

		go func() {
			defer close(results)
			defer close(errCh)

			_ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			for {
				select {
				case query := <-in:
					q := query.(datastore.Query)
					func() {
						q := f.ApplyQuery(m.client.Collection(m.
							collectionName), q)

						iter := t.Documents(q)
						defer iter.Stop()

						for {
							ref, err := iter.Next()

							if err == iterator.Done {
								return
							}

							var usr User
							if err = ref.DataTo(&usr); err != nil {
								errCh <- err
								cancel()
								return
							}

							results <- &usr
						}
					}()
				case <-_ctx.Done():
					return
				}
			}
		}()

		return results, errCh
	}
}

func (m *FirestoreManager) find(ctx context.Context,
	t *firestore.Transaction, users chan<- *User, errs chan<- error,
	queries ...datastore.Query) {
	_ctx, cancel := context.WithCancel(ctx)
	producer := make(chan datastore.Query)
	go func() {
		for _, query := range queries {
			producer <- query
		}
	}()

	p := pipeline.NewPipelineWithProducer[datastore.Query, *User](producer)
	if _, err := p.AddStageStreamWorker(5, 5,
		m.createStreamWorker(t)); err != nil {
		errs <- err
		cancel()
		return
	}

	out, err := p.Produces()
	if err != nil {
		errs <- err
		cancel()
		return
	}

	go func() {
		_errCh := p.Start(_ctx)
		for {
			select {
			case e := <-_errCh:
				errs <- e
				cancel()
			case u := <-out:
				users <- u
			case <-_ctx.Done():
				return
			}
		}
	}()
}

// Add adds a user to the database of users.
//
// Please note that a user is considered a duplicate if any of the following
// already exist on a different user: Email, PhoneNumber and Username.The Add
// method will return ErrUserDuplicate in this case.
func (m *FirestoreManager) Add(ctx context.Context,
	usr *User) (<-chan *User, <-chan error) {
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
		userCh <- usr
	}()

	return userCh, errCh
}

// Find finds the user based on the given query criterion.
func (m *FirestoreManager) Find(ctx context.Context,
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
func (m *FirestoreManager) Update(ctx context.Context,
	user *User) (<-chan *User, <-chan error) {
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
		userCh <- user
	}()

	return userCh, errCh
}
