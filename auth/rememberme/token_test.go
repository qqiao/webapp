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

package rememberme_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/qqiao/webapp/auth/rememberme"
)

var tm rememberme.TokenManager

func setUp() {
	client, err := firestore.NewClient(context.Background(), "test-project")
	if err != nil {
		log.Fatalf("Unable to initialize firebase client. Error: %v", err)
	}

	tm = rememberme.NewTokenManager(client, "TestCollection")
}

func TestMain(m *testing.M) {
	setUp()
	os.Exit(m.Run())
}

func TestTokenManager_Save(t *testing.T) {
	identifier, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Unable to create UUID for token. Error: %v", err)
	}

	newToken := rememberme.Token{
		Username:   "test_user",
		Identifier: identifier.String(),
	}

	tokenCh, errCh := tm.Save(context.Background(), newToken)
	select {
	case err := <-errCh:
		t.Errorf("Error saving token: %v", err)
	case <-tokenCh:
	}

	// Once a token is saved, subsequent validation calls should succeed
	tokenCh, errCh = tm.Validate(context.Background(), newToken)
	select {
	case err = <-errCh:
		t.Errorf("Error validating token: %v", err)
	case <-tokenCh:
	}
}

func TestTokenManager_SaveToken(t *testing.T) {
	identifier, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Unable to create UUID for token. Error: %v", err)
	}

	newToken := rememberme.Token{
		Username:   "test_user",
		Identifier: identifier.String(),
	}

	errCh := tm.SaveToken(context.Background(), newToken)
	select {
	case err = <-errCh:
		if err != nil {
			t.Errorf("Error saving token: %v", err)
		}
	}

	// Once a token is saved, subsequent validation calls should succeed
	errCh = tm.ValidateToken(context.Background(), newToken)
	select {
	case err = <-errCh:
		if err != nil {
			t.Errorf("Error validating token: %v", err)
		}
	}
}

func TestTokenManager_ValidateToken(t *testing.T) {
	identifier, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Unable to create UUID for token. Error: %v", err)
	}

	newToken := rememberme.Token{
		Username:   "test_user",
		Identifier: identifier.String(),
	}

	errCh := tm.SaveToken(context.Background(), newToken)

	select {
	case err = <-errCh:
		if err != nil {
			t.Errorf("Error saving token: %v", err)
		}
	}

	// Non-existent users shouldn't have valid tokens
	errCh = tm.ValidateToken(context.Background(), rememberme.Token{
		Username: "non_existent",
	})
	select {
	case err = <-errCh:
		if err != nil && err != rememberme.ErrTokenInvalid {
			t.Errorf("Error validating token: %v", err)
		}
	}

	// Validation should fail for valid users with non-existent identifiers
	errCh = tm.ValidateToken(context.Background(), rememberme.Token{
		Username:   "test_user",
		Identifier: "non_existent",
	})
	select {
	case err = <-errCh:
		if err != nil && err != rememberme.ErrTokenInvalid {
			t.Errorf("Error validating token: %v", err)
		}
	}
}

func TestTokenManager_RevokeToken(t *testing.T) {
	identifier, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Unable to create UUID for token. Error: %v", err)
	}

	newToken := rememberme.Token{
		Username:   "test_user",
		Identifier: identifier.String(),
	}

	errCh := tm.SaveToken(context.Background(), newToken)
	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Error saving token: %v", err)
		}
	}

	// token should validate at this point
	errCh = tm.ValidateToken(context.Background(), newToken)
	select {
	case err = <-errCh:
		if err != nil {
			t.Errorf("Error validating token: %v", err)
		}
	}

	errCh = tm.RevokeToken(context.Background(), newToken)
	select {
	case err = <-errCh:
		if err != nil {
			t.Errorf("Error revoking token: %v", err)
		}
	}

	// Validation should fail after the token has been revoked
	errCh = tm.ValidateToken(context.Background(), newToken)
	select {
	case err = <-errCh:
		if err != rememberme.ErrTokenInvalid {
			t.Errorf("Revoked token should not validat: %v", err)
		}
	}
}

func TestTokenManager_PurgeTokens(t *testing.T) {
	oldIdentifier, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Unable to create UUID for token. Error: %v", err)
	}
	oldToken := rememberme.Token{
		Username:   "test_user",
		Identifier: oldIdentifier.String(),
	}
	errCh := tm.SaveToken(context.Background(), oldToken)
	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Error saving token: %v", err)
		}
	}

	time.Sleep(5 * time.Second)
	cutoff := time.Now()
	time.Sleep(5 * time.Second)

	newIdentifier, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Unable to create UUID for token. Error: %v", err)
	}
	newToken := rememberme.Token{
		Username:   "test_user",
		Identifier: newIdentifier.String(),
	}
	errCh = tm.SaveToken(context.Background(), newToken)
	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Error saving token: %v", err)
		}
	}

	errCh = tm.PurgeTokens(context.Background(), "test_user", cutoff)
	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Error purging token: %v", err)
		}
	}

	// oldToken should now have been purged and validation should fail
	errCh = tm.ValidateToken(context.Background(), oldToken)
	select {
	case err = <-errCh:
		if err != nil && err != rememberme.ErrTokenInvalid {
			t.Errorf("Error validating token: %v", err)
		}
	}

	// newToken should still validate
	errCh = tm.ValidateToken(context.Background(), newToken)
	select {
	case err = <-errCh:
		if err != nil {
			t.Errorf("New token should still validate, but failed: %v", err)
		}
	}

	// Purging non-existent tokens shouldn't matter
	errCh = tm.PurgeTokens(context.Background(), "non-existent", time.Now())
	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Error purging token: %v", err)
		}
	}

	// after purging non-existent tokens,
	// the ones that do fail should still faile,
	// and ones that do work should still work'

	// oldToken should now have been purged and validation should fail
	errCh = tm.ValidateToken(context.Background(), oldToken)
	select {
	case err = <-errCh:
		if err != nil && err != rememberme.ErrTokenInvalid {
			t.Errorf("Error validating token: %v", err)
		}
	}

	// newToken should still validate
	errCh = tm.ValidateToken(context.Background(), newToken)
	select {
	case err = <-errCh:
		if err != nil {
			t.Errorf("New token should still validate, but failed: %v", err)
		}
	}
}
