package endpoints

import (
	"context"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleOAuth2Callback(t *testing.T) {
	assertion := assert.New(t)

	// see oauth2_test in golang.org/x/oauth2
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Header.Get("Authorization"), "Basic Q0xJRU5UX0lEJTNGJTNGOkNMSUVOVF9TRUNSRVQlM0YlM0Y="; got != want {
			t.Errorf("Authorization header = %q; want %q", got, want)
		}

		w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		_, _ = w.Write([]byte("access_token=90d64460d14870c08c81352a05dedd3465940a7c&scope=user&token_type=bearer"))
	}))
	defer ts.Close()

	request, err := http.NewRequest("GET", "someUrl", nil)
	if err != nil {
		panic(err)
	}

	// see oauth2_test in golang.org/x/oauth2
	config := oauth2.Config{
		ClientID:     "CLIENT_ID??",
		ClientSecret: "CLIENT_SECRET??",
		RedirectURL:  "RedirectMe",
		Endpoint: oauth2.Endpoint{
			AuthURL:  ts.URL + "/auth",
			TokenURL: ts.URL + "/token",
		},
	}
	ctx := context.Background()

	auth := auth{ctx: &ctx, oauth2Config: OAuth2Config{config: &config}}
	responseRecorder := httptest.NewRecorder()

	auth.handleOAuth2Callback(responseRecorder, request)

	assertion.Equal(http.StatusOK, responseRecorder.Code)
	assertion.Equal("", responseRecorder.Body.String())
}

func TestHandleLogin(t *testing.T) {
	assertion := assert.New(t)

	request, err := http.NewRequest("GET", "someUrl", nil)
	if err != nil {
		panic(err)
	}

	config := oauth2.Config{ClientID: "TestClient", RedirectURL: "RedirectMe"}
	auth := auth{oauth2Config: OAuth2Config{config: &config}}
	responseRecorder := httptest.NewRecorder()

	auth.handleLogin(responseRecorder, request)

	assertion.Equal(http.StatusFound, responseRecorder.Code)
	assertion.Equal(".?client_id=TestClient&redirect_uri=RedirectMe&response_type=code&state=main", responseRecorder.Header().Get("Location"))
}
