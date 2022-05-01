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
	"github.com/qqiao/webapp/v2/auth/rememberme"
)

var tm *rememberme.FirestoreTokenManager

func setUp() {
	client, err := firestore.NewClient(context.Background(), "test-project")
	if err != nil {
		log.Fatalf("Unable to initialize firebase client. Error: %v", err)
	}

	tm = rememberme.NewFirestoreTokenManager(client, "TestTokenCollection")
}

func TestMain(m *testing.M) {
	setUp()
	os.Exit(m.Run())
}

func TestFirestoreTokenManager_Add(t *testing.T) {
	identifier := uuid.NewString()

	token := rememberme.Token{
		Username:   "test_user",
		Identifier: identifier,
	}

	t.Run("Addition should succeed", func(t *testing.T) {
		tokenCh, errCh := tm.Add(context.Background(), token)
		select {
		case err := <-errCh:
			t.Errorf("Error saving token: %v", err)
		case <-tokenCh:
		}
	})

	t.Run("Added token should validate", func(t *testing.T) {
		tokenCh, errCh := tm.Validate(context.Background(), token)
		select {
		case err := <-errCh:
			t.Errorf("Error validating token: %v", err)

		case <-tokenCh:
		}
	})

	t.Run("Adding the same token should produce ErrTokenDuplicate",
		func(t *testing.T) {
			tokenCh, errCh := tm.Add(context.Background(), token)
			select {
			case err := <-errCh:
				if err != rememberme.ErrTokenDuplicate {
					t.Errorf("Error saving token: %v", err)
				}
			case <-tokenCh:
				t.Errorf("Saving the same token again should error out")
			}
		})
}

func TestFirestoreTokenManager_Delete(t *testing.T) {
	identifier := uuid.NewString()
	token := rememberme.Token{
		Username:   "test_user",
		Identifier: identifier,
	}

	t.Run("Initial addition", func(t *testing.T) {
		tokenCh, errCh := tm.Add(context.Background(), token)
		select {
		case err := <-errCh:
			t.Errorf("Error saving token: %v", err)

		case <-tokenCh:
		}

		// token should validate at this point
		tokenCh, errCh = tm.Validate(context.Background(), token)
		select {
		case err := <-errCh:
			t.Errorf("Error validating token: %v", err)

		case <-tokenCh:
		}
	})

	t.Run("Deletion should succeed", func(t *testing.T) {
		errCh := tm.Delete(context.Background(), token)
		select {
		case err := <-errCh:
			if err != nil {
				t.Errorf("Error deleting token: %v", err)
			}
		}
	})

	t.Run("Validation should fail on deleted token", func(t *testing.T) {
		tokenCh, errCh := tm.Validate(context.Background(), token)
		select {
		case err := <-errCh:
			if err != rememberme.ErrTokenInvalid {
				t.Errorf("Revoked token should not validat: %v", err)
			}

		case <-tokenCh:
			t.Error("Validation should have failed for revoked token")
		}
	})
}

func TestFirestoreTokenManager_Purge(t *testing.T) {
	oldIdentifier1 := uuid.NewString()
	oldToken1 := rememberme.Token{
		Username:   "test_user",
		Identifier: oldIdentifier1,
	}

	tokenCh, errCh := tm.Add(context.Background(), oldToken1)
	select {
	case err := <-errCh:
		t.Errorf("Error saving token: %v", err)
	case <-tokenCh:
	}

	oldIdentifier2 := uuid.NewString()
	oldToken2 := rememberme.Token{
		Username:   "test_user",
		Identifier: oldIdentifier2,
	}

	tokenCh, errCh = tm.Add(context.Background(), oldToken2)
	select {
	case err := <-errCh:
		t.Errorf("Error saving token: %v", err)
	case <-tokenCh:
	}

	time.Sleep(5 * time.Second)
	cutoff := time.Now()
	time.Sleep(5 * time.Second)

	newIdentifier := uuid.NewString()
	newToken := rememberme.Token{
		Username:   "test_user",
		Identifier: newIdentifier,
	}

	tokenCh, errCh = tm.Add(context.Background(), newToken)
	select {
	case err := <-errCh:
		t.Errorf("Error saving token: %v", err)
	case <-tokenCh:
	}

	errCh = tm.Purge(context.Background(), "test_user", cutoff)
	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Error purging token: %v", err)
		}
	}

	// oldToken1 should now have been purged and validation should fail
	tokenCh, errCh = tm.Validate(context.Background(), oldToken1)
	select {
	case err := <-errCh:
		if err != rememberme.ErrTokenInvalid {
			t.Errorf("Error validating token: %v", err)
		}
	case <-tokenCh:
		t.Errorf("Shouldn't get a token as validation should have failed")
	}

	// oldToken1 should now have been purged and validation should fail
	tokenCh, errCh = tm.Validate(context.Background(), oldToken2)
	select {
	case err := <-errCh:
		if err != rememberme.ErrTokenInvalid {
			t.Errorf("Error validating token: %v", err)
		}
	case <-tokenCh:
		t.Errorf("Shouldn't get a token as validation should have failed")
	}

	// newToken should still validate
	tokenCh, errCh = tm.Validate(context.Background(), newToken)
	select {
	case err := <-errCh:
		t.Errorf("New token should still validate, but failed: %v", err)

	case <-tokenCh:
	}

	// Purging non-existent tokens shouldn't matter
	errCh = tm.Purge(context.Background(), "non-existent", time.Now())
	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Error purging token: %v", err)
		}
	}

	// after purging non-existent tokens,
	// the ones that do fail should still fail,
	// and ones that do work should still work'

	// oldToken should now have been purged and validation should fail
	tokenCh, errCh = tm.Validate(context.Background(), oldToken1)
	select {
	case err := <-errCh:
		if err != rememberme.ErrTokenInvalid {
			t.Errorf("Error validating token: %v", err)
		}
	case <-tokenCh:
		t.Error("Validation should have failed")
	}

	// newToken should still validate
	tokenCh, errCh = tm.Validate(context.Background(), newToken)
	select {
	case err := <-errCh:
		t.Errorf("New token should still validate, but failed: %v", err)

	case <-tokenCh:
	}
}

func TestFirestoreTokenManager_Revoke(t *testing.T) {
	identifier := uuid.NewString()

	newToken := rememberme.Token{
		Username:   "test_user",
		Identifier: identifier,
	}

	tokenCh, errCh := tm.Add(context.Background(), newToken)
	select {
	case err := <-errCh:
		t.Errorf("Error saving token: %v", err)

	case <-tokenCh:
	}

	// token should validate at this point
	tokenCh, errCh = tm.Validate(context.Background(), newToken)
	select {
	case err := <-errCh:
		t.Errorf("Error validating token: %v", err)

	case <-tokenCh:
	}

	tokenCh, errCh = tm.Revoke(context.Background(), newToken)
	select {
	case err := <-errCh:
		t.Errorf("Error revoking token: %v", err)

	case tok := <-tokenCh:
		if !tok.Revoked {
			t.Error("Revoked flag should now be set as true")
		}
	}

	// Validation should fail after the token has been revoked
	tokenCh, errCh = tm.Validate(context.Background(), newToken)
	select {
	case err := <-errCh:
		if err != rememberme.ErrTokenInvalid {
			t.Errorf("Revoked token should not validat: %v", err)
		}

	case <-tokenCh:
		t.Error("Validation should have failed for revoked token")
	}
}

func TestFirestoreTokenManager_Validate(t *testing.T) {
	identifier := uuid.NewString()

	token := rememberme.Token{
		Username:   "test_user",
		Identifier: identifier,
	}

	tokenCh, errCh := tm.Add(context.Background(), token)

	select {
	case err := <-errCh:
		t.Errorf("Error saving token: %v", err)

	case t := <-tokenCh:
		token = *t
	}

	// Non-existent users shouldn't have valid tokens
	tokenCh, errCh = tm.Validate(context.Background(), rememberme.Token{
		Username: "non_existent",
	})
	select {
	case err := <-errCh:
		if err != rememberme.ErrTokenInvalid {
			t.Errorf("Error validating token: %v", err)
		}
	case <-tokenCh:
		t.Error("Non existent user shouldn't have any valid tokens")
	}

	// Validation should fail for valid users with non-existent identifiers
	tokenCh, errCh = tm.Validate(context.Background(), rememberme.Token{
		Username:   "test_user",
		Identifier: "non_existent",
	})
	select {
	case err := <-errCh:
		if err != rememberme.ErrTokenInvalid {
			t.Errorf("Error validating token: %v", err)
		}

	case <-tokenCh:
		t.Error("Validation should fail for bad ID")
	}

	time.Sleep(5 * time.Second)

	// Validation should succeed, and last used should get updated
	tokenCh, errCh = tm.Validate(context.Background(), rememberme.Token{
		Username:   "test_user",
		Identifier: identifier,
	})
	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Validation should succeed, but failed: %v", err)
		}

	case newT := <-tokenCh:
		if token.LastUsed >= newT.LastUsed {
			t.Error("Last used should have been updated after validation")
		}
	}
}
