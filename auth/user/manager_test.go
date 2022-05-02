package user_test

import (
	"context"
	"testing"

	"github.com/qqiao/webapp/v2/auth/user"
	"github.com/qqiao/webapp/v2/datastore"
)

type methodTestFunc func(user.Manager) func(*testing.T)

var managers = make(map[string]user.Manager)
var tests = map[string]methodTestFunc{
	"Add":    testAdd,
	"Find":   testFind,
	"Update": testUpdate,
}

func TestManagers(t *testing.T) {
	for managerName, manager := range managers {
		t.Run(managerName, func(t *testing.T) {
			for methodName, test := range tests {
				t.Run(methodName, test(manager))
			}
		})
	}
}

func testAdd(m user.Manager) func(*testing.T) {
	return func(t *testing.T) {
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
}

func testFind(m user.Manager) func(*testing.T) {
	return func(t *testing.T) {
		const username = "test_find"
		usr := user.NewUser().WithUsername(username).WithPassword("123")

		// First time adding the user should succeed
		userCh, errCh := m.Add(context.Background(), usr)
		select {
		case err := <-errCh:
			t.Errorf("Error adding user: %v", err)
		case <-userCh:
		}

		t.Run("Should be able to retrieve existing user", func(t *testing.T) {
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

		t.Run("Should be able to retreve with laxer criterion", func(t *testing.T) {
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

				for range u {
					count++
				}
				if count > 1 {
					t.Errorf("Should have only found 1 user, got: %d", count)
				}
			}
		})

		t.Run("Should not find phantoms", func(t *testing.T) {
			foundCh, errCh := m.Find(context.Background(), datastore.Query{
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
		})
	}
}

func testUpdate(m user.Manager) func(*testing.T) {
	return func(t *testing.T) {
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
}
