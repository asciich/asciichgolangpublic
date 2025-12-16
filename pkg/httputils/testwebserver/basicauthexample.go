package testwebserver

import (
	"encoding/json"
	"net/http"

	"github.com/asciich/asciichgolangpublic/pkg/httputils/basicauth"
	"github.com/asciich/asciichgolangpublic/pkg/randomgenerator"
)

// This is an example implementation using basic auth.
// It's used by the testwebserver.
//
// For an implementation of a BasicAuth function to protect your endpoints use the basichauth.BasicAuth() function.
type BasicAuthExample struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Get a new basic auth example.
// It is automatically initialized with an example user and a random generated password.
func NewBasicAuthExample() *BasicAuthExample {
	pw, err := randomgenerator.GetRandomString(10)
	if err != nil {
		panic(err)
	}

	return &BasicAuthExample{
		Username: "exampleuser",
		Password: pw,
	}
}

func (b *BasicAuthExample) IndexHtml(w http.ResponseWriter, r *http.Request) {
	content := `<html>
<h1>BasicAuthExample</h1>
This is a basic auth example:
<ul>
	<li><a href="credentials.json">credentials.json</a> will return the currently used credentials in JSON format.</li>
	<li><a href="protected.txt">protected.txt</a> is a basic auth protected endpoint you can only reach successfully with the username '` + b.Username + `' and password '` + b.Password + `' (or take the credentials from <a href="credentials.json">credentials.json</a>).</li>
</ul>
</html>`

	w.Write([]byte(content))
}

func (b *BasicAuthExample) CredentialsJson(w http.ResponseWriter, r *http.Request) {
	content, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		panic(err)
	}

	w.Write(content)
}

func (b *BasicAuthExample) Protected(w http.ResponseWriter, r *http.Request) {
	protectedEndpoint := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is a basic auth protected message."))
	}
	basicauth.BasicAuthSingleCredentials(protectedEndpoint, b.Username, b.Password)(w, r)
}

func (b *BasicAuthExample) GetUsername(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(b.Username))
}

func (b *BasicAuthExample) GetPassword(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(b.Password))
}
