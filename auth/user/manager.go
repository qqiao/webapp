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

// Manager is responsible for all user related operations
type Manager interface {
	// Add adds a user to the database of users.
	//
	// Add will return ErrUserDuplicate if the user already exists  in the
	// datastore.
	Add(ctx context.Context, user User) <-chan error

	Find(ctx context.Context,
		query datastore.Query) (<-chan *User, <-chan error)
}
