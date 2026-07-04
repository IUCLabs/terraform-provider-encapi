package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientSendsBearerToken(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		_, _ = w.Write([]byte(`{"name":"roles::base"}`))
	}))
	defer srv.Close()

	c := newAPIClient(srv.URL, "tok")
	var out roleAPI
	if err := c.put(context.Background(), "/api/v1/roles/roles::base", roleAPI{Name: "roles::base"}, &out); err != nil {
		t.Fatal(err)
	}
	if gotAuth != "Bearer tok" {
		t.Errorf("auth = %q, want Bearer tok", gotAuth)
	}
}

func TestClientNoTokenNoHeader(t *testing.T) {
	var hadAuth bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, hadAuth = r.Header["Authorization"]
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := newAPIClient(srv.URL, "")
	if err := c.get(context.Background(), "/api/v1/roles/x", &roleAPI{}); err != nil {
		t.Fatal(err)
	}
	if hadAuth {
		t.Error("no Authorization header expected when token empty")
	}
}

func TestClientNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	err := newAPIClient(srv.URL, "").get(context.Background(), "/api/v1/nodes/ghost", &nodeAPI{})
	if !isNotFound(err) {
		t.Errorf("err = %v, want notFound", err)
	}
}

func TestClientAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"role and environment are required"}`))
	}))
	defer srv.Close()

	err := newAPIClient(srv.URL, "t").put(context.Background(), "/api/v1/nodes/h", nodeAPI{}, &nodeAPI{})
	if err == nil {
		t.Fatal("expected api error")
	}
	if isNotFound(err) {
		t.Errorf("400 should not be a notFound error, got %v", err)
	}
}
