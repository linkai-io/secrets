package secrets_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/linkai-io/secrets"
)

const expected = "password123"

func TestSecretsCacheLocal(t *testing.T) {
	env := "local"
	key := fmt.Sprintf("/some/%s/password", env)

	s := secrets.NewSecretsCache(env, "")

	if err := s.SetSecureParameter(key, expected); err != nil {
		t.Fatalf("error setting password: %v\n", err)
	}

	returned, err := s.Password()
	if err != nil {
		t.Fatalf("error reading password: %v\n", err)
	}

	if returned != expected {
		t.Fatalf("expected %q got %q", expected, returned)
	}
}

func TestSecretsCacheAWS(t *testing.T) {
	if os.Getenv("INFRA_TESTS") == "" {
		t.Skip("skipping infrastructure tests")
	}

	env := "dev"
	region := "us-east-1"
	key := fmt.Sprintf("/some/%s/password", env)

	s := secrets.NewSecretsCache(env, region)
	if err := s.SetSecureParameter(key, expected); err != nil {
		t.Fatalf("error setting password: %v\n", err)
	}

	returned, err := s.Password()
	if err != nil {
		t.Fatalf("error reading password: %v\n", err)
	}

	if returned != expected {
		t.Fatalf("expected %q got %q", expected, returned)
	}
}
