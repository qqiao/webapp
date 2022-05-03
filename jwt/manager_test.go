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
	"testing"
	"time"

	"github.com/qqiao/webapp/v2/jwt"
)

type testCase struct {
	manager  jwt.Manager
	token    string
	expected string
	dat      []string
}

var testCases []testCase

func testSigning(t *testing.T, manager jwt.Manager, dat string) {
	claims := jwt.NewClaims().WithDat(dat).WithExpiry(time.Unix(253402271999, 0))
	tok, errCh := manager.SignCustom(claims)

	select {
	case err := <-errCh:
		t.Errorf("Failed to create token: %v", err)
	case token := <-tok:
		decodedCh, errCh := manager.ParseCustom(token)

		select {
		case err := <-errCh:
			t.Errorf("Failed to validate token: %v", err)

		case decoded := <-decodedCh:
			if dat != decoded.Dat {
				t.Errorf("Did not get back input.\nInput: %q\nGot: %q",
					dat, decoded.Dat)
			}
		}
	}
}

// func FuzzSigning(f *testing.F) {
// 	f.Add("1")

// 	f.Fuzz(testMatch)
// }

func testParseCustom(t *testing.T, manager jwt.Manager, token string, expected string) {
	gotCh, errCh := manager.ParseCustom(token)

	select {
	case err := <-errCh:
		t.Errorf("Error while validating token: %v", err)

	case got := <-gotCh:
		if got.Dat != expected {
			t.Errorf("Expected: %s. Got: %s", expected, got.Dat)
		}
	}
}

func testRepeatableParseCustoms(t *testing.T, manager jwt.Manager, token string, expected string) {
	for i := 0; i < 10; i++ {
		gotCh, errCh := manager.ParseCustom(token)

		select {
		case err := <-errCh:
			t.Errorf("Error while validating token: %v", err)

		case got := <-gotCh:
			if got.Dat != expected {
				t.Errorf("Expected: %s. Got: %s", expected, got.Dat)
			}
		}
	}
}

func TestManagers(t *testing.T) {
	for _, tc := range testCases {
		alg := tc.manager.Alg()
		t.Run(fmt.Sprintf("%s should parse correctly", alg),
			func(t *testing.T) {
				testParseCustom(t, tc.manager, tc.token, tc.expected)
			})

		t.Run(fmt.Sprintf("%s should perform repeatable parses", alg),
			func(t *testing.T) {
				testRepeatableParseCustoms(t, tc.manager, tc.token, tc.expected)
			})

		for _, dat := range tc.dat {
			t.Run(fmt.Sprintf("%s should sign correctly", alg),
				func(t *testing.T) {
					testSigning(t, tc.manager, dat)
				})
		}
	}
}
