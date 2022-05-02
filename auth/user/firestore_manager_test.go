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

package user_test

import (
	"context"
	"log"
	"os"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/qqiao/webapp/v2/auth/user"
	"github.com/qqiao/webapp/v2/datastore"
)

var m *user.FirestoreManager

func setUp() {
	client, err := firestore.NewClient(context.Background(), "test-project")
	if err != nil {
		log.Fatalf("Unable to initialize firebase client. Error: %v", err)
	}

	m = user.NewFirestoreManager(client, "TestUserCollection")
}

func TestMain(m *testing.M) {
	setUp()
	os.Exit(m.Run())
}

func TestFirestoreManager_Add(t *testing.T) {
	const username = "test_add_user"
	usr := user.NewUser().WithUsername(username).WithPassword("123")

	t.Run("Initial Add should succeed", func(t *testing.T) {
		userCh, errCh := m.Add(context.Background(), usr)
		select {
		case err := <-errCh:
			t.Errorf("Error adding user: %v", err)
		case u := <-userCh:
			if u.UID == "" {
				t.Error("Newly added user should have a generated UID")
			}
		}
	})

	t.Run("Adding again should get ErrUserDuplicate", func(t *testing.T) {
		userCh, errCh := m.Add(context.Background(), usr)
		select {
		case err := <-errCh:
			if err != user.ErrUserDuplicate {
				t.Errorf("Expecting ErrUserDuplicate, got: %v", err)
			}
		case <-userCh:
			t.Error("Adding the same user again should result in ErrUserDuplicate")
		}
	})

	t.Run("Added users should be retrievable", func(t *testing.T) {
		foundCh, errCh := m.Find(context.Background(), datastore.Query{
			Filters: []datastore.Filter{
				{
					Path:     "Username",
					Operator: "==",
					Value:    username,
				},
				{
					Path:     "Password",
					Operator: "==",
					Value:    "123",
				},
			},
		})
		select {
		case err := <-errCh:
			if err != nil {
				t.Errorf("Error when finding user: %v", err)
			}
		case u := <-foundCh:
			count := 0

			for range u {
				count++
			}
			if count > 1 {
				t.Errorf("Should have only found 1 user, got: %d", count)
			}
		}
	})
}

func TestFirestoreManager_Find(t *testing.T) {
	const username = "test_find"
	usr := user.NewUser().WithUsername(username).WithPassword("123")

	// First time adding the user should succeed
	userCh, errCh := m.Add(context.Background(), usr)
	select {
	case err := <-errCh:
		t.Errorf("Error adding user: %v", err)
	case <-userCh:
	}

	// We should be able to retrieve the user
	foundCh, errCh := m.Find(context.Background(), datastore.Query{
		Filters: []datastore.Filter{
			{
				Path:     "Username",
				Operator: "==",
				Value:    username,
			},
			{
				Path:     "Password",
				Operator: "==",
				Value:    "123",
			},
		},
	})
	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Error when finding user: %v", err)
		}
	case u := <-foundCh:
		count := 0

		for range u {
			count++
		}
		if count > 1 {
			t.Errorf("Should have only found 1 user, got: %d", count)
		}
	}

	// We should be able to retrieve the user with less restrictive criteria
	foundCh, errCh = m.Find(context.Background(), datastore.Query{
		Filters: []datastore.Filter{
			{
				Path:     "Username",
				Operator: "==",
				Value:    username,
			},
		},
	})
	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Error when finding user: %v", err)
		}
	case u := <-foundCh:
		count := 0

		for range u {
			count++
		}
		if count > 1 {
			t.Errorf("Should have only found 1 user, got: %d", count)
		}
	}

	// We shouldn't be able to find non-existent users
	foundCh, errCh = m.Find(context.Background(), datastore.Query{
		Filters: []datastore.Filter{
			{
				Path:     "Username",
				Operator: "==",
				Value:    "non-existent",
			},
		},
	})
	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Error when finding user: %v", err)
		}
	case u := <-foundCh:
		count := 0

		for range u {
			count++
		}
		if count > 0 {
			t.Errorf("Should have found 0 users, got: %d", count)
		}
	}
}

func TestFirestoreManager_Update(t *testing.T) {
	const username = "test_update"
	usr := user.NewUser().WithUsername(username).WithPassword("123").
		WithSuspended(false)

	// First time adding the user should succeed
	userCh, errCh := m.Add(context.Background(), usr)
	select {
	case err := <-errCh:
		t.Errorf("Error adding user: %v", err)
	case <-userCh:
	}

	// Then we are going to update the record
	usr.Suspended = true
	userCh, errCh = m.Update(context.Background(), usr)
	select {
	case err := <-errCh:
		t.Errorf("Error updating user: %v", err)
	case <-userCh:
	}

	// Let's retrieve the user and compare
	foundCh, errCh := m.Find(context.Background(), datastore.Query{
		Filters: []datastore.Filter{
			{
				Path:     "Username",
				Operator: "==",
				Value:    username,
			},
		},
	})
	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Error when finding user: %v", err)
		}
	case u := <-foundCh:
		count := 0
		var lastFound *user.User
		for foundUser := range u {
			count++
			lastFound = foundUser
		}
		if count > 1 {
			t.Errorf("Should have only found 1 user, got: %d", count)
		}

		if lastFound.Username != usr.Username ||
			lastFound.Password != usr.Password ||
			lastFound.Suspended != usr.Suspended {
			t.Errorf("Updated failed.\nExpected: %v\nGot: %v", usr, *lastFound)
		}
	}
}
