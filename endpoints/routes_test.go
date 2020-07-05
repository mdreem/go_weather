package endpoints

import (
	"context"
	"github.com/coreos/go-oidc"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type HttpHandlerTestData struct {
}

func (httpHandlerTestData HttpHandlerTestData) ServeHTTP(http.ResponseWriter, *http.Request) {
}

func TestAuthenticationHandlerAccessTokenMissing(t *testing.T) {
	assertion := assert.New(t)
	request := http.Request{Header: make(http.Header)}

	auth := auth{}
	responseRecorder := httptest.NewRecorder()
	auth.authenticationHandler(HttpHandlerTestData{}, responseRecorder, &request)

	assertion.Equal(http.StatusUnauthorized, responseRecorder.Code)
}

func TestAuthenticationHandlerAccessTokenWrongFormat(t *testing.T) {
	assertion := assert.New(t)
	request := http.Request{Header: make(http.Header)}
	request.Header.Add("Authorization", "OnlyOnePart")

	auth := auth{}
	responseRecorder := httptest.NewRecorder()
	auth.authenticationHandler(HttpHandlerTestData{}, responseRecorder, &request)

	assertion.Equal(http.StatusBadRequest, responseRecorder.Code)
}

type testVerifier struct{}

func (t *testVerifier) VerifySignature(_ context.Context, _ string) ([]byte, error) {
	return nil, errors.New("verification failed")
}

func TestAuthenticationHandlerAccessTokenCannotBeVerified(t *testing.T) {
	assertion := assert.New(t)

	request := http.Request{Header: make(http.Header)}
	request.Header.Add("Authorization", "Bearer: invalidToken")

	keySet := &testVerifier{}
	config := oidc.Config{
		SkipClientIDCheck: true,
		SkipExpiryCheck:   true,
		SkipIssuerCheck:   true,
	}
	failingVerifier := oidc.NewVerifier("someIssuer", keySet, &config)

	ctx := context.Background()

	auth := auth{ctx: &ctx, verifier: failingVerifier}
	responseRecorder := httptest.NewRecorder()
	auth.authenticationHandler(HttpHandlerTestData{}, responseRecorder, &request)

	assertion.Equal(http.StatusUnauthorized, responseRecorder.Code)
}
