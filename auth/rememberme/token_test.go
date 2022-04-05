package rememberme_test

import (
	"context"
	"log"
	"os"
	"testing"

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
	case err := <-errCh:
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
		if err != rememberme.ErrTokenInvalid {
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
		if err != rememberme.ErrTokenInvalid {
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
