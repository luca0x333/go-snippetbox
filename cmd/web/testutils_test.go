package main

import (
	"github.com/golangcollege/sessions"
	"github.com/luca0x333/go-snippetbox/pkg/models/mock"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"
)

// newTestApplication helper returns an instance of the application struct containing mocked dependencies.
func newTestApplication(t *testing.T) *application {
	// Create an instance of the template cache.
	templateCache, err := newTemplateCache("./../../ui/html/")
	if err != nil {
		t.Fatal(err)
	}

	// Create a session manager instance.
	session := sessions.New([]byte("z6Nah+pPonzHbI*+9Pk8qNWhTzbpa@ge"))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	// Initialize the dependencies using the mocks for the loggers and database models.
	return &application{
		errorLog:      log.New(ioutil.Discard, "", 0),
		infoLog:       log.New(ioutil.Discard, "", 0),
		session:       session,
		snippets:      &mock.SnippetModel{},
		templateCache: templateCache,
		users:         &mock.UserModel{},
	}
}

// Custom testServer type which embeds a httptest.Server instance.
type testServer struct {
	*httptest.Server
}

// newTestServer helper initializes and returns a new instance of our custom testServer.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	// Initialize a new cookie jar.
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add the cookie jar to the client so that response cookies are stored and then sent with the requests.
	ts.Client().Jar = jar

	// Disable redirect-following for the client.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

// get() method makes a GET request to a given url path on the testServer and returns the response status code,
// headers and body.
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}
