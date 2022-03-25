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
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	jw "github.com/qqiao/webapp/jwt"
)

const token = `eyJhbGciOiJQUzUxMiIsInR5cCI6IkpXVCJ9.eyJkYXQiOiIxIiwiZXhwIjoyNTM0MDIyNzE5OTl9.acGc_2bhg9ELH0YCumW8BBrsI7nNXUw2CJWMOJVRXbYCGY3BvKBWkNKFuy-q_zRZ8RDlN1qI0oKokakHxk_94Gg8x7ttJbVg5-dysL3hhS0E5eZGpX40ujSSqW5s1bctBjOjAFU9weR7DKSqznglMgUL6_K11I2F8ZG3aTTtc8wFMN3D1wplqiw3RhbLbsyFJx8p2ZEokIzofNP7SIUcmKyXuVx9_me9BRdfTH8mwJ4miSfyW8Aq9vASGWYb8TDuTlPi4yGTrzzjvzdG8OLyfkoK4oaK_6uW2ZzAwkXFjMLiy1RuRkj36aH5IOGSoBdS8ns32wfeOu8mTOzn_dOa2ztIQD_iwX5z-3kcx_v1emAzvsPro7p6yPjE75Z5qU0rw7EgHYvCigg96hLs1ghNRHFN4Xx5ahMl4dqDJPA0L6EQsj80mqfDgAJ7285jYpZs28X7Ij19fqRoVw-fvsj_zcEI4WJnhapY9pbiOwbh8EUxtltgW3IiPzKLohgAF8JZ6rnnJJqOWi9TGbknfeLh6cXkohWMTlk8q6uu9g25SLdravvCUReFvIkJYIukO2y8wDPTlB9gOR9uQcdTKn-Wr6G43GS05hhappKjotAqxuvlaMEdaVHh_Qr1fLcy7erMd69irR7dbMsfZ5BriEyWE9OTAr8Ano7qoXMZlqt-37Q`

var j jw.Manager

func setUp() {
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(`
-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`))

	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(`
-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAwe3SUOlXW3TRxOs+CfJb
9xABVCSW9LdjRKAvJvcAvbR5nVVXfv078fVL+9a/mr+V2FzPXi/QRW7QeFEBT9gO
pljYxRWH8T+6hA2UETrDaYGsjEcjl1YridwlH5elVsn0tvcCE2B9lKYhAYwsMx3q
mcaCUXonC21aa+uncGdrxCFkxh1uosdCy81eaXMU5hyuDYRsddVcyG9XaZSFKmcP
m2IFG1rwEDrl8AXjAo+h+u/7ekSzYCwIckp6VbJ/FyD3p7OG8RtTIycbxceVaZcY
AfjS2volBV+pjpsa5+LvU4AuYu4lBOcfOyzLkIQUb6EVTnSVw4H7efahLkAfv3dW
lbzu3gmYlFo5cd/tCqdXfm+xWYOxGOw6ko7aS3C5W/zYBapyTfNnsHLpz447TvKz
ekYzcyu+bueFEUpCta3O6pi57HdI95vNEcLz7hS1yzQme+zBY7liuQWQ/r0dVAnG
ksZnaHodfr/vGmeer0vahkpDgzzk2ucbyX1n6YM6nnlFcCcRWkYs6bArz2xswhuu
X38Ffkxe5i0ICG3OJks8de/1qb3fcFv5Tlf6U9c/giZsrpjxDuK05QR+LOSHd6Yy
TEZPQCoOkZTAV0DmzZGu3LAsSZFJDu/Yp345CyWm1sY5RvqM4pyRcrgTRIdvoBPu
rxYt4E6zaoD1Ix9YXh2bKocCAwEAAQ==
-----END PUBLIC KEY-----
`))

	j = jw.NewManager(publicKey, privateKey)
}

func TestMain(m *testing.M) {
	setUp()
	os.Exit(m.Run())
}

func testMatch(t *testing.T, input string) {
	tm := time.Unix(253402271999, 0)
	tok, errCh := j.CreateCustom(input, &tm)

	select {
	case err := <-errCh:
		t.Errorf("Failed to create token: %v", err)
	case token := <-tok:
		decodedCh, errCh := j.ValidateCustom(token)

		select {
		case err := <-errCh:
			t.Errorf("Failed to validate token: %v", err)

		case decoded := <-decodedCh:
			if input != decoded {
				t.Errorf("Did not get back input.\nInput: %q\nGot: %q",
					input, decoded)
			}
		}
	}

}

// func FuzzSigning(f *testing.F) {
// 	f.Add("1")

// 	f.Fuzz(testMatch)
// }

func TestSigning(t *testing.T) {
	testMatch(t, "1")
	testMatch(t, "中文")
}

func TestValidation(t *testing.T) {
	expected := "1"
	gotCh, errCh := j.ValidateCustom(token)

	select {
	case err := <-errCh:
		t.Errorf("Error while validating token: %v", err)

	case got := <-gotCh:
		if got != expected {
			t.Errorf("Expected: %s. Got: %s", expected, got)
		}
	}
}

func TestMultipleValidations(t *testing.T) {
	for i := 0; i < 10; i++ {
		expected := "1"
		gotCh, errCh := j.ValidateCustom(token)

		select {
		case err := <-errCh:
			t.Errorf("Error while validating token: %v", err)

		case got := <-gotCh:
			if got != expected {
				t.Errorf("Expected: %s. Got: %s", expected, got)
			}
		}
	}
}
