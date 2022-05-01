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

package jwt_test

import (
	"fmt"
	"time"

	"github.com/qqiao/webapp/v2/jwt"
)

func ExampleClaims_WithDat() {
	claims := jwt.NewClaims().WithDat("123")

	fmt.Println(claims.Dat)

	// Output: 123
}

func ExampleClaims_WithExpiry() {
	now := time.Unix(0, 0).Add(1 * time.Hour)
	claims := jwt.NewClaims().WithExpiry(now)
	fmt.Printf("%d", claims.ExpiresAt)

	// Output: 3600
}
