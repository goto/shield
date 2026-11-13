package basic_auth

import (
	"crypto/sha1"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	goauth "github.com/abbot/go-http-auth"

	"github.com/goto/salt/log"
	"github.com/goto/shield/core/rule"
	"github.com/goto/shield/internal/proxy/middleware"
	"github.com/goto/shield/pkg/httputil"
	"golang.org/x/crypto/bcrypt"
)

func TestBcryptAuthenticationAndLegacyRejection(t *testing.T) {
	password := "synthetic-test-only"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	sha := sha1.Sum([]byte(password)) // synthetic fixture for rejecting legacy SHA-1
	legacySHA := "{SHA}" + base64.StdEncoding.EncodeToString(sha[:])
	legacyMD5 := string(goauth.MD5Crypt([]byte(password), []byte("testsalt"), []byte("$apr1$")))
	for _, legacy := range []string{legacySHA, legacyMD5} {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.SetBasicAuth("test-user", password)
		original := goauth.NewBasicAuthenticator("shield", func(string, string) string { return legacy })
		if original.CheckAuth(r) != "test-user" {
			t.Fatal("legacy test fixture must be valid")
		}
	}
	for _, tc := range []struct {
		name, hash, password string
		allowed              bool
	}{
		{"bcrypt-success", string(hash), password, true},
		{"bcrypt-wrong-password", string(hash), "synthetic-wrong", false},
		{"sha1-rejected", legacySHA, password, false},
		{"md5-rejected", legacyMD5, password, false},
		{"malformed-rejected", "invalid", password, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				if r.Header.Get(httputil.HeaderXUser) != "test-user" {
					t.Fatal("authenticated identity not forwarded")
				}
				w.WriteHeader(http.StatusNoContent)
			})
			auth := New(log.NewNoop(), next)
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.SetBasicAuth("test-user", tc.password)
			middleware.EnrichRule(r, &rule.Rule{Middlewares: rule.MiddlewareSpecs{{Name: "basic_auth", Config: map[string]interface{}{"users": []map[string]interface{}{{"user": "test-user", "password": tc.hash}}}}}})
			w := httptest.NewRecorder()
			auth.ServeHTTP(w, r)
			if called != tc.allowed {
				t.Fatal("unexpected downstream forwarding")
			}
			if !tc.allowed && w.Code != http.StatusUnauthorized {
				t.Fatalf("status %d", w.Code)
			}
		})
	}
}

func TestRoutesWithoutBasicAuthKeepForwarding(t *testing.T) {
	called := false
	auth := New(log.NewNoop(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true; w.WriteHeader(http.StatusNoContent) }))
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	middleware.EnrichRule(r, &rule.Rule{})
	w := httptest.NewRecorder()
	auth.ServeHTTP(w, r)
	if !called || w.Code != http.StatusNoContent {
		t.Fatal("unrelated route was blocked")
	}
}
