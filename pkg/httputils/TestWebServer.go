package httputils

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httputilsimplementationindependend"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httputilsinterfaces"
	"github.com/asciich/asciichgolangpublic/tlsutils/x509utils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

// A simple webserver mostly used for testing.
type TestWebServer struct {
	webServerWaitGroup *sync.WaitGroup
	port               int
	mux                *http.ServeMux
	server             *http.Server
	tlsConfig          *tls.Config
}

func GenerateCertAndKeyForTestWebserver(ctx context.Context) (certAndKeyPair *x509utils.X509CertKeyPair, err error) {
	return x509utils.CreateSelfSignedCertificate(
		ctx,
		&x509utils.X509CreateCertificateOptions{
			Organization:   "localorg",
			CommonName:     "localhost",
			CountryName:    "CH",
			Locality:       "Zurich",
			AdditionalSans: []string{"localhost"},
		},
	)
}

func GetTlsTestWebServer(ctx context.Context, port int) (webServer httputilsinterfaces.Server, err error) {
	toReturn := NewTestWebServer()

	err = toReturn.SetPort(port)
	if err != nil {
		return nil, err
	}

	certAndKey, err := GenerateCertAndKeyForTestWebserver(ctx)
	if err != nil {
		return nil, err
	}

	err = toReturn.SetTlsCertAndKey(ctx, certAndKey)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

func (t *TestWebServer) GetTlsCert() (cert *x509.Certificate, err error) {
	if t.tlsConfig == nil {
		return
	}

	if len(t.tlsConfig.Certificates) != 1 {
		return
	}

	tlsCert := t.tlsConfig.Certificates[0]

	return x509utils.TlsCertToX509Cert(&tlsCert)
}

func GetTestWebServer(port int) (webServer httputilsinterfaces.Server, err error) {
	toReturn := NewTestWebServer()

	err = toReturn.SetPort(port)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

func NewTestWebServer() (t *TestWebServer) {
	return new(TestWebServer)
}

func (t *TestWebServer) GetMux() (mux *http.ServeMux, err error) {
	if t.mux == nil {
		return nil, tracederrors.TracedErrorf("mux not set")
	}

	return t.mux, nil
}

func (t *TestWebServer) GetPort() (port int, err error) {
	if t.port <= 0 {
		return -1, tracederrors.TracedError("port not set")
	}

	return t.port, nil
}

func (t *TestWebServer) GetServer() (server *http.Server, err error) {
	if t.server == nil {
		return nil, tracederrors.TracedErrorf("server not set")
	}

	return t.server, nil
}

func (t *TestWebServer) GetWebServerWaitGroup() (webServerWaitGroup *sync.WaitGroup, err error) {
	if t.webServerWaitGroup == nil {
		return nil, tracederrors.TracedErrorf("webServerWaitGroup not set")
	}

	return t.webServerWaitGroup, nil
}

func (t *TestWebServer) SetMux(mux *http.ServeMux) (err error) {
	if mux == nil {
		return tracederrors.TracedErrorf("mux is nil")
	}

	t.mux = mux

	return nil
}

func (t *TestWebServer) SetPort(port int) (err error) {
	if port <= 0 {
		return tracederrors.TracedErrorf("Invalid value '%d' for port", port)
	}

	t.port = port

	return nil
}

func (t *TestWebServer) SetServer(server *http.Server) (err error) {
	if server == nil {
		return tracederrors.TracedErrorf("server is nil")
	}

	t.server = server

	return nil
}

func (t *TestWebServer) SetWebServerWaitGroup(webServerWaitGroup *sync.WaitGroup) (err error) {
	if webServerWaitGroup == nil {
		return tracederrors.TracedErrorf("webServerWaitGroup is nil")
	}

	t.webServerWaitGroup = webServerWaitGroup

	return nil
}

func (t *TestWebServer) StartInBackground(ctx context.Context) (err error) {
	port, err := t.GetPort()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Start testWebServer in background on port %d started.", port)

	if t.webServerWaitGroup == nil {
		t.webServerWaitGroup = new(sync.WaitGroup)
	} else {
		return tracederrors.TracedError(httputilsimplementationindependend.ErrWebServerAlreadyRunning)
	}

	t.mux = http.NewServeMux()
	t.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r == nil {
			logging.LogWarn("r is nil")
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		if r.URL == nil {
			logging.LogWarn("r.URL is nil")
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		io.WriteString(w, "TestWebServer main page\n")
	})

	t.mux.HandleFunc("/hello_world.txt", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello world\n")
	})

	t.mux.HandleFunc("/example1.yaml", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "---\nhello: world\n")
	})

	t.mux.HandleFunc("/example1.json", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"hello": "world"}`)
	})

	t.server = &http.Server{
		Addr:      ":" + strconv.Itoa(port),
		Handler:   t.mux,
		TLSConfig: t.tlsConfig,
	}

	t.webServerWaitGroup.Add(1)
	go func() {
		defer t.webServerWaitGroup.Done()

		if t.tlsConfig == nil {
			if err := t.server.ListenAndServe(); err != http.ErrServerClosed {
				log.Fatalf("ListenAndServe(): %v", err)
			}
		} else {
			if err := t.server.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
				log.Fatalf("ListenAndServeTLS(): %v", err)
			}
		}
	}()

	time.Sleep(1 * time.Second)

	logging.LogInfoByCtxf(ctx, "Start testWebServer in background on port %d finished.", port)

	return nil
}

func (t *TestWebServer) SetTlsCertAndKey(ctx context.Context, certAndKey *x509utils.X509CertKeyPair) (err error) {
	if certAndKey == nil {
		return tracederrors.TracedErrorNil("certAndKey")
	}

	err = certAndKey.CheckKeyMatchingCert()
	if err != nil {
		return err
	}

	tlsCert := tls.Certificate{
		Certificate: [][]byte{certAndKey.Cert.Raw},
		PrivateKey:  certAndKey.Key,
	}

	t.tlsConfig = &tls.Config{Certificates: []tls.Certificate{tlsCert}}

	return nil
}

func (t *TestWebServer) Stop(ctx context.Context) (err error) {
	logging.LogInfoByCtx(ctx, "Stop TestWebServer started.")

	if t.webServerWaitGroup == nil {
		logging.LogInfoByCtxf(ctx, "TestWebServer already stopped")
		return nil
	}

	if t.server == nil {
		return tracederrors.TracedError("Unexpected t.server == nil")
	}

	err = t.server.Shutdown(context.TODO())
	if err != nil {
		return tracederrors.TracedErrorf(
			"Shutdown TestWebServer failed: '%w'",
			err,
		)
	}

	t.webServerWaitGroup.Wait()
	t.webServerWaitGroup = nil

	logging.LogInfoByCtx(ctx, "Stop TestWebServer finished.")

	return nil
}
