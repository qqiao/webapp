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

// User represents a user record stored in the underlying datastore.
type User struct {
	DisplayName string `json:"displayName,omitempty"`
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	PhotoURL    string `json:"photoUrl,omitempty"`
	Suspended   bool   `json:"-"`
	UID         string `json:"uid,omitempty" firestore:"-"`
	Username    string `json:"username,omitempty"`
}

// NewUser returns a new User instance with all fields being blank.
func NewUser() *User {
	return &User{}
}

// WithDisplayName sets the DisplayName of the user.
func (u *User) WithDisplayName(displayName string) *User {
	u.DisplayName = displayName
	return u
}

// WithEmail sets the Email of the user.
func (u *User) WithEmail(email string) *User {
	u.Email = email
	return u
}

// WithPassword sets the Password of the user.
func (u *User) WithPassword(password string) *User {
	u.Password = password
	return u
}

// WithPhoneNumber sets the PhoneNumber of the user.
func (u *User) WithPhoneNumber(phoneNumber string) *User {
	u.PhoneNumber = phoneNumber
	return u
}

// WithPhotoURL sets the PhotoURL of the user.
func (u *User) WithPhotoURL(photoURL string) *User {
	u.PhotoURL = photoURL
	return u
}

// WithSuspended sets the Suspended of the user.
func (u *User) WithSuspended(suspended bool) *User {
	u.Suspended = suspended
	return u
}

// WithUID sets the UID of the user.
func (u *User) WithUID(uid string) *User {
	u.UID = uid
	return u
}

// WithUsername sets the Username of the user.
func (u *User) WithUsername(username string) *User {
	u.Username = username
	return u
}
