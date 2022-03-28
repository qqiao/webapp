package jwt_test

import (
	"fmt"
	"time"

	"github.com/qqiao/webapp/jwt"
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
