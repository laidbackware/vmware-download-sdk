package myvmware_test

import (
	"os"
	"testing"

	"github.com/hoegaarden/myvmware"
)

func TestE2E(t *testing.T) {
	user, pass := mustEnv(t, "VMW_USER"), mustEnv(t, "VMW_PASS")

	client, err := myvmware.Login(user, pass)
	if err != nil {
		t.Errorf("expected error not to occur, got %q", err)
	}

	t.Logf("Logged in")

	currentUser, err := client.CurrentUser()
	if err != nil {
		t.Errorf("expected error not to occur, got: %q", err)
	}

	mustNotBeEmpty(t, "currentUser response", currentUser)

	t.Logf("got current user: %#v", currentUser)

	accountInfo, err := client.AccountInfo()
	if err != nil {
		t.Errorf("expected account info not to response with an error, got: %q", err)
	}

	mustNotBeEmpty(t, "accountInfo response", accountInfo)

	t.Logf("got account info: %#v", accountInfo)

	products, err := client.Products()
	if err != nil {
		t.Errorf("expected error not to occur, got: %q", err)
	}

	mustNotBeEmpty(t, "products response", products)

	t.Logf("got products: %#v", products)
}

func mustEnv(t *testing.T, k string) string {
	t.Helper()

	if v, ok := os.LookupEnv(k); ok {
		return v
	}

	t.Fatalf("expected environment variable %q", k)
	return ""
}

// TODO: Just for now, until we have proper types for responses
func mustNotBeEmpty(t *testing.T, name string, raw interface{}) {
	t.Helper()

	data, ok := raw.(map[string]interface{})
	if !ok {
		t.Errorf("expected %q to be a map, got %T", name, raw)
	}

	if len(data) == 0 {
		t.Errorf("expected %q not to be empty", name)
	}
}
