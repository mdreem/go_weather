package endpoints

import (
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"testing"
)

type OAuth2TestConfig struct {
	exchangeError error
	token         *oauth2.Token
}

func (config OAuth2TestConfig) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return config.token, config.exchangeError
}

func (config OAuth2TestConfig) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	return ""
}

func (config OAuth2TestConfig) GetClientID() string {
	return ""
}

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

func TestHandleOAuth2CallbackExchangeFailed(t *testing.T) {
	assertion := assert.New(t)

	request, err := http.NewRequest("GET", "someUrl", nil)
	if err != nil {
		panic(err)
	}
	responseRecorder := httptest.NewRecorder()

	ctx := context.Background()
	auth := auth{
		oauth2Config: OAuth2TestConfig{
			exchangeError: errors.New("exchange error"),
		},
		ctx: &ctx,
	}

	auth.handleOAuth2Callback(responseRecorder, request)

	assertion.Equal(500, responseRecorder.Code)
}

func TestHandleOAuth2CallbackTokenMissing(t *testing.T) {
	assertion := assert.New(t)

	request, err := http.NewRequest("GET", "someUrl", nil)
	if err != nil {
		panic(err)
	}
	responseRecorder := httptest.NewRecorder()

	ctx := context.Background()
	tokenMap := make(map[string]interface{})
	token := oauth2.Token{}
	tokenWithExtraData := token.WithExtra(tokenMap)
	auth := auth{
		oauth2Config: OAuth2TestConfig{
			token: tokenWithExtraData,
		},
		ctx: &ctx,
	}

	auth.handleOAuth2Callback(responseRecorder, request)

	assertion.Equal(500, responseRecorder.Code)
}

func TestHandleOAuth2CallbackWriteResponseFailed(t *testing.T) {
	assertion := assert.New(t)

	request, err := http.NewRequest("GET", "someUrl", nil)
	if err != nil {
		panic(err)
	}
	responseRecorder := httptest.NewRecorder()

	ctx := context.Background()
	tokenMap := make(map[string]interface{})
	tokenMap["id_token"] = "someToken"
	token := oauth2.Token{}
	tokenWithExtraData := token.WithExtra(tokenMap)
	auth := auth{
		oauth2Config: OAuth2TestConfig{
			token: tokenWithExtraData,
		},
		ctx: &ctx,
	}

	auth.handleOAuth2Callback(responseRecorder, request)

	assertion.Equal(http.StatusOK, responseRecorder.Code)
	assertion.Equal("someToken", responseRecorder.Body.String())
}
