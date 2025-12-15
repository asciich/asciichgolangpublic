package httpoptions

import (
	"encoding/base64"
	"fmt"
)

type BasicAuth struct {
	Username string
	Password string
}

// Returns the value of the Authorization header.
// It's int eh format of "Basic <base64_encoded_credentials>"
func (b *BasicAuth) AuthorizationValue() string {
	authString := fmt.Sprintf("%s:%s", b.Username, b.Password)
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(authString))
	return "Basic " + encodedAuth
}
