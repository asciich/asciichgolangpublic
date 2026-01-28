package publicips

import (
	"context"
	"encoding/json"

	"github.com/asciich/asciichgolangpublic/pkg/httputils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const GET_PUBLIC_IP_URL = "https://asciich.ch/what_is_my_ip.php"

// Get the own public IP address.
// Usfull when behind a NAS to get the public address from inside the natted network.
func GetPublicIp(ctx context.Context) (string, error) {
	logging.LogInfoByCtxf(ctx, "Get public IP address started.")

	response, err := httputils.SendRequest(ctx, &httpoptions.RequestOptions{
		Url: GET_PUBLIC_IP_URL,
	})
	if err != nil {
		return "", err
	}

	decoded := &struct {
		ClientIp string `json:"client_ip"`
	}{}

	data, err := response.GetBodyAsString()
	if err != nil {
		return "", err
	}

	err = json.Unmarshal([]byte(data), decoded)
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to get public IP using %s : %w", GET_PUBLIC_IP_URL, err)
	}

	ip := decoded.ClientIp

	logging.LogInfoByCtxf(ctx, "Get public IP address finished. Public IP is '%s'.", ip)

	return ip, nil
}
