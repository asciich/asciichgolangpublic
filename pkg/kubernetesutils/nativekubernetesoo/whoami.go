package nativekubernetesoo

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubeconfigutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	authenticationv1 "k8s.io/api/authentication/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

// WhoAmI returns the identity used to authenticate against the cluster.
//
// Strategy:
//  1. Ask the API server using the SelfSubjectReview API. This is the most
//     reliable way since the server reports the identity.
//  2. If the API is not available (cluster < 1.27, or disabled on 1.27) the server
//     answers with NotFound. In this case the identity is derived from the
//     credentials of the rest config (see whoAmIByRestConfig).
//
// Errors other than NotFound (e.g. cluster not reachable, unauthorized) are
// returned as they are and do NOT trigger the fallback.
//
// The cluster name is optional (e.g. for the default cluster) and only used
// for log and error messages.
func (n *NativeKubernetesCluster) WhoAmI(ctx context.Context) (*kubernetesgeneric.UserInfo, error) {
	clusterName := n.displayName()

	clientset, err := n.GetClientSet()
	if err != nil {
		return nil, err
	}

	var username string

	response, err := clientset.AuthenticationV1().SelfSubjectReviews().Create(ctx, &authenticationv1.SelfSubjectReview{}, metav1.CreateOptions{})
	if err == nil {
		username = response.Status.UserInfo.Username
	} else {
		if !errors.IsNotFound(err) {
			return nil, tracederrors.TracedErrorf("SelfSubjectReview for cluster '%s' failed: %w", clusterName, err)
		}

		logging.LogInfoByCtxf(
			ctx,
			"Whoami: SelfSubjectReview API not available in cluster '%s'. Fall back to rest config based detection.",
			clusterName,
		)

		username, err = n.whoAmIByRestConfig(ctx)
		if err != nil {
			return nil, err
		}
	}

	logging.LogInfoByCtxf(ctx, "Whoami: Kube config uses user '%s' to log in to cluster '%s'.", username, clusterName)

	return &kubernetesgeneric.UserInfo{
		Username: username,
	}, nil
}

// displayName returns the cluster name for log and error messages.
// The default cluster has no explicit name, so a placeholder is returned instead of an error.
func (n *NativeKubernetesCluster) displayName() string {
	if n.name == "" {
		return "<default cluster>"
	}

	return n.name
}

// resolveKubeContext returns the kubectl context of the cluster.
// For the default cluster (no name) the current-context of the kubeconfig is used.
func (n *NativeKubernetesCluster) resolveKubeContext(ctx context.Context) (string, error) {
	if n.name != "" {
		return n.GetKubectlContext(ctx)
	}

	rawConfig, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		return "", tracederrors.TracedErrorf("Unable to load kubeconfig: %w", err)
	}

	if rawConfig.CurrentContext == "" {
		return "", tracederrors.TracedError("No current-context set in kubeconfig and no cluster name given.")
	}

	return rawConfig.CurrentContext, nil
}

// whoAmIByRestConfig determines the user from the credentials in the rest config.
//
// Order:
//  1. CN of the client certificate (default for kind / kubeadm clusters).
//     This is the username the API server sees.
//  2. "sub" claim of a bearer token if it is a JWT
//     (e.g. "system:serviceaccount:<namespace>:<name>" for in-cluster auth).
//  3. Username of basic authentication.
//  4. User name configured in the kubeconfig (e.g. exec plugins, OIDC).
//     For the default cluster the current-context of the kubeconfig is used.
//
// This is only the identity the client *claims*, not one confirmed by the server.
func (n *NativeKubernetesCluster) whoAmIByRestConfig(ctx context.Context) (string, error) {
	config, err := n.GetConfig()
	if err != nil {
		return "", err
	}

	// 1. Client certificate
	certPEM := config.TLSClientConfig.CertData
	if len(certPEM) == 0 && config.TLSClientConfig.CertFile != "" {
		certPEM, err = os.ReadFile(config.TLSClientConfig.CertFile)
		if err != nil {
			return "", tracederrors.TracedErrorf(
				"Unable to read client certificate file '%s': %w",
				config.TLSClientConfig.CertFile,
				err,
			)
		}
	}

	if len(certPEM) > 0 {
		cn, err := commonNameFromPEMCertificate(certPEM)
		if err != nil {
			return "", tracederrors.TracedErrorf("Unable to get CN from client certificate: %w", err)
		}

		if cn != "" {
			return cn, nil
		}
	}

	// 2. Bearer token (service account tokens are JWTs)
	if config.BearerToken != "" {
		subject, err := subjectFromJWT(config.BearerToken)
		if err != nil {
			// Not a JWT (opaque token), continue with the next option.
			logging.LogInfoByCtxf(ctx, "Whoami: Bearer token is not a JWT with a subject: %v", err)
		} else if subject != "" {
			return subject, nil
		}
	}

	// 3. Basic auth
	if config.Username != "" {
		return config.Username, nil
	}

	// 4. User name from the kubeconfig
	kubeContext, err := n.resolveKubeContext(ctx)
	if err != nil {
		return "", err
	}

	logging.LogInfoByCtxf(ctx, "Whoami: No certificate or token identity found. Using the kubeconfig user name of context '%s'.", kubeContext)

	return kubeconfigutils.GetUserNameByContextName(ctx, kubeContext)
}

// commonNameFromPEMCertificate returns the subject CN of the first certificate
// in the given PEM data (the leaf certificate for a client cert chain).
func commonNameFromPEMCertificate(pemBytes []byte) (string, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return "", fmt.Errorf("no PEM block found in client certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", err
	}

	return cert.Subject.CommonName, nil
}

// subjectFromJWT returns the "sub" claim of a JWT without verifying the signature.
// Verification is not needed: the value is only used to name the identity
// the client is configured with.
func subjectFromJWT(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("token is not a JWT")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	var claims struct {
		Subject string `json:"sub"`
	}

	err = json.Unmarshal(payload, &claims)
	if err != nil {
		return "", err
	}

	return claims.Subject, nil
}
