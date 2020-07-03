package endpoints

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleLogin(t *testing.T) {
	assertion := assert.New(t)

	request, err := http.NewRequest("GET", "someUrl", nil)
	if err != nil {
		panic(err)
	}

	config := oauth2.Config{ClientID: "TestClient", RedirectURL: "RedirectMe"}
	auth := auth{oauth2Config: &config}
	responseRecorder := httptest.NewRecorder()

	auth.handleLogin(responseRecorder, request)

	assertion.Equal(http.StatusFound, responseRecorder.Code)
	assertion.Equal(".?client_id=TestClient&redirect_uri=RedirectMe&response_type=code&state=main", responseRecorder.Header().Get("Location"))
}
