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

package rememberme

import "errors"

// Errors
var (
	ErrTokenDuplicate = errors.New("token duplicate")
	ErrTokenNotFound  = errors.New("token not found")
	ErrTokenInvalid   = errors.New("token invalid")
)

// Token represents a rememberme token stored.
type Token struct {
	Username   string
	Identifier string
	Revoked    bool
	UserAgent  string
	Created    int64
	LastUsed   int64
}
