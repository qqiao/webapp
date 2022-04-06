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

package inhouse

import (
	"github.com/qqiao/webapp/auth/user"
)

// User represents a user to be stored in the datastore.
//
// Deprecated: please use the auth/user package instead.
type User = user.User

// Errors.
//
// Deprecated: please use the auth/user package instead.
var (
	ErrUserExists   = user.ErrUserDuplicate
	ErrUserNotFound = user.ErrUserNotFound
)

// UserManager is responsible for all user related operations
//
// Deprecated: please use the auth/user package instead.
type UserManager = user.FirebaseManager

// NewUserManager creates a new UserManager with the given firestore client
// and collection name.
//
// Deprecated: please use the auth/user package instead.
var NewUserManager = user.NewFirebaseManager
