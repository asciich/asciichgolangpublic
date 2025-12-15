package basicauth_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/basicauth"
)

func Test_BasicAuthWithSingleUsernameAndPassword(t *testing.T) {
	// create a new http serveMux
	mux := http.NewServeMux()

	// We use a single username and password to protect our endpoint.
	const username = "testuser"
	const password = "testpassword"

	// This function handles the request after a successful authentication.
	// Define your actual logic and payload to return in this function:
	endpoint := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is a basic auth protected message."))
	}

	// Add the endpoint "/protected.txt" which directly vcalls basciauth.BasicAuthSingleCredentials to do the
	mux.HandleFunc("/protected.txt", func(w http.ResponseWriter, r *http.Request) {
		// The basicauth.BasicAuthSingleCredentials is used to protect our protectedEndpoint with basic auth:
		basicauth.BasicAuthSingleCredentials(endpoint, username, password)(w, r)
	})

	// Create the webserver with our endpoint:
	const port = 12345
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: mux,
	}

	// start our example webserver in the background:
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			panic(fmt.Sprintf("ListenAndServe(): %v", err))
		}
	}()
	// stop the example server on exit.
	defer server.Shutdown(context.Background())

	// give webserver time to start in the background
	time.Sleep(time.Millisecond * 200)

	// This is the url used for testing
	url := fmt.Sprintf("http://localhost:%d/protected.txt", port)

	// Try unauthenticated request:
	response, err := http.Get(url)
	require.NoError(t, err)
	defer response.Body.Close()
	// The unauthenticated requests with an unauthorized status code:
	assert.EqualValues(t, http.StatusUnauthorized, response.StatusCode)

	// Try with the wrong credentials:
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", "wronguser", "wrongpass")))
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	request.Header.Set("Authorization", fmt.Sprintf("Basic %s", encodedAuth))
	response2, err := new(http.Client).Do(request)
	require.NoError(t, err)
	defer response2.Body.Close()
	// The unauthenticated requests with an unauthorized status code:
	assert.EqualValues(t, http.StatusUnauthorized, response2.StatusCode)

	// Try with the correct credentials:
	encodedAuth2 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
	request2, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	request2.Header.Set("Authorization", fmt.Sprintf("Basic %s", encodedAuth2))
	response3, err := new(http.Client).Do(request2)
	require.NoError(t, err)
	defer response3.Body.Close()
	// With the correct username and password we are able to do a successful request.
	assert.EqualValues(t, http.StatusOK, response3.StatusCode)
	// And therefore we receive to message defined in the handler function at the beginning of this test:
	content, err := io.ReadAll(response3.Body)
	require.NoError(t,err)
	require.EqualValues(t, "This is a basic auth protected message.", string(content))
}
