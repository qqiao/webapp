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
	"errors"

	"github.com/qqiao/webapp/datastore"
)

// Errors.
var (
	ErrUserDuplicate = errors.New("duplicate user")
	ErrUserNotFound  = errors.New("user not found")
)

// Manager is responsible for all user related operations.
//
// Depending on how users are stored and queried, there could be multiple
// implementations of the Manager interface.
type Manager interface {
	// Add adds a user to the database of users.
	//
	// Please note that a user is considered a duplicate if any of the following
	// already exist on a different user: Email, PhoneNumber and Username.The Add
	// method will return ErrUserDuplicate in this case.
	Add(ctx context.Context, user *User) (<-chan *User, <-chan error)

	// Find finds the user based on the given query criterion
	Find(ctx context.Context,
		query datastore.Query) (<-chan (<-chan *User), <-chan error)

	// Update updates the given user record.
	//
	// Update will return ErrUserNotFound if the user cannot be found in the
	// underlying datastore
	Update(ctx context.Context, user *User) (<-chan *User, <-chan error)
}
