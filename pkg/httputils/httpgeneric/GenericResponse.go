package httpgeneric

import (
	"context"
	"crypto/x509"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/fileformats/jsonutils"
	"github.com/asciich/asciichgolangpublic/pkg/fileformats/yamlutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// This is the generic response type.
// It can also be seen as the default response to use.
type GenericResponse struct {
	body               []byte
	statusCode         int
	serverCertificates []*x509.Certificate
}

func NewGenericResponse() (g *GenericResponse) {
	return new(GenericResponse)
}

func (g *GenericResponse) RunYqQueryAgainstBody(query string) (result string, err error) {
	if query == "" {
		return "", tracederrors.TracedErrorEmptyString("query")
	}

	body, err := g.GetBodyAsString()
	if err != nil {
		return "", err
	}

	return yamlutils.RunYqQueryAginstYamlStringAsString(body, query)
}

func (g *GenericResponse) RunJqQueryAgainstBody(query string) (result string, err error) {
	if query == "" {
		return "", tracederrors.TracedErrorEmptyString("query")
	}

	body, err := g.GetBodyAsString()
	if err != nil {
		return "", err
	}

	return jsonutils.RunJqAgainstJsonStringAsString(body, query)
}

func (g *GenericResponse) GetBody() (body []byte, err error) {
	if g.body == nil {
		return nil, tracederrors.TracedErrorf("body not set")
	}

	return g.body, nil
}

func (g *GenericResponse) GetBodyAsBytes() ([]byte, error) {
	return g.GetBody()
}

func (g *GenericResponse) GetBodyAsString() (body string, err error) {
	bodyBytes, err := g.GetBody()
	if err != nil {
		return "", err
	}

	return string(bodyBytes), err
}

func (g *GenericResponse) GetStatusCode() (statusCode int, err error) {
	if g.statusCode <= 0 {
		return -1, tracederrors.TracedError("statusCode not set")
	}

	return g.statusCode, nil
}

func (g *GenericResponse) CheckStatusCode(expectedStatusCodes []int) error {
	for _, toCheck := range expectedStatusCodes {
		if g.IsStatusCode(toCheck) {
			return nil
		}
	}

	return tracederrors.TracedErrorf("%w: %d does not match expected status codes %v. Response body is:\n%s", ErrUnexpectedStatusCode, g.statusCode, expectedStatusCodes, g.body)
}

func (g *GenericResponse) IsStatusCode(expectedStatusCode int) bool {
	statusCode, err := g.GetStatusCode()
	if err != nil {
		return false
	}

	return statusCode == expectedStatusCode
}

func (g *GenericResponse) IsStatusCode200Ok() bool {
	return g.IsStatusCode(http.StatusOK)
}

func (g *GenericResponse) SetBody(body []byte) (err error) {
	if body == nil {
		return tracederrors.TracedErrorf("body is nil")
	}

	g.body = body

	return nil
}

func (g *GenericResponse) SetStatusCode(statusCode int) (err error) {
	if statusCode <= 0 {
		return tracederrors.TracedErrorf("Invalid value '%d' for statusCode", statusCode)
	}

	g.statusCode = statusCode

	return nil
}

func (g *GenericResponse) GetServerEndEntitiyCertificate(ctx context.Context) (*x509.Certificate, error) {
	if g.serverCertificates == nil || len(g.serverCertificates) == 0 {
		return nil, tracederrors.TracedError("no certificates collected")
	}

	return g.serverCertificates[0], nil
}

func (g *GenericResponse) GetServerCertificateChain(ctx context.Context) ([]*x509.Certificate, error) {
	if g.serverCertificates == nil {
		return nil, tracederrors.TracedError("no certificates collected")
	}

	return g.serverCertificates, nil
}

func (g *GenericResponse) LogCertInfo(ctx context.Context) error {
	if g.serverCertificates == nil || len(g.serverCertificates) == 0 {
		return tracederrors.TracedError("no certificates collected")
	}

	logging.LogInfoByCtxf(ctx, "Server certificate chain contains %d certificate(s):", len(g.serverCertificates))

	for i, cert := range g.serverCertificates {
		if i == 0 {
			logging.LogInfoByCtxf(ctx, "Certificate %d (End Entity/Leaf):", i+1)
		} else if i == len(g.serverCertificates)-1 {
			logging.LogInfoByCtxf(ctx, "Certificate %d (Root CA):", i+1)
		} else {
			logging.LogInfoByCtxf(ctx, "Certificate %d (Intermediate CA):", i+1)
		}

		logging.LogInfoByCtxf(ctx, "  Subject: %s", cert.Subject.CommonName)
		logging.LogInfoByCtxf(ctx, "  Issuer: %s", cert.Issuer.CommonName)
		logging.LogInfoByCtxf(ctx, "  Serial Number: %s", formatSerialNumber(cert.SerialNumber))
		logging.LogInfoByCtxf(ctx, "  Valid From: %s", cert.NotBefore.Format("2006-01-02 15:04:05 UTC"))
		logging.LogInfoByCtxf(ctx, "  Valid Until: %s", cert.NotAfter.Format("2006-01-02 15:04:05 UTC"))

		if len(cert.DNSNames) > 0 {
			logging.LogInfoByCtxf(ctx, "  DNS Names: %v", cert.DNSNames)
		}
		if len(cert.IPAddresses) > 0 {
			ips := make([]string, len(cert.IPAddresses))
			for j, ip := range cert.IPAddresses {
				ips[j] = ip.String()
			}
			logging.LogInfoByCtxf(ctx, "  IP Addresses: %v", ips)
		}
	}

	return nil
}

// formatSerialNumber formats the certificate serial number with colon separators
func formatSerialNumber(serial *big.Int) string {
	hexBytes := serial.Bytes()
	parts := make([]string, len(hexBytes))
	for i, b := range hexBytes {
		parts[i] = fmt.Sprintf("%02X", b)
	}
	return strings.Join(parts, ":")
}

func (g *GenericResponse) SetServerCertificates(certificates []*x509.Certificate) error {
	if certificates == nil {
		return tracederrors.TracedErrorNil("certificates")
	}

	g.serverCertificates = certificates

	return nil
}
