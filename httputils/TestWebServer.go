package httputils

import (
	"context"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

// A simple webserver mostly used for testing.
type TestWebServer struct {
	webServerWaitGroup *sync.WaitGroup
	port               int
	mux                *http.ServeMux
	server             *http.Server
}

/* TODO implement
func GetTlsTestWebServer(port int, verbose bool) (webServer Server, err error) {
	toReturn, err := GetTestWebServer(port)
	if err != nil {
		return nil, err
	}

	certFilePath, err := tempfiles.CreateEmptyTemporaryFile(verbose)
	if err != nil {
		return nil, err
	}

	keyFilePath, err := tempfiles.CreateEmptyTemporaryFileAndGetPath(verbose)
	if err != nil {
		return nil, err
	}
	defer files.MustDeleteFilesByPath(
		verbose,
		certFilePath,
		keyFilePath,
	)

	err := x509utils.CreateSelfSignedCertificate(
		&x509utils.X509CreateCertificateOptions{
			CommonName:                "localhost",
			CountryName:               "CH",
			Locality:                  "Zurich",
			AdditionalSans:            []string{"localhost"},
			Verbose:                   true,
			KeyOutputFilePath:         keyFilePath,
			CertificateOutputFilePath: certFilePath,
		},
	)

	err = toReturn.SetTlsCertAndKey(certFilePath, keyFilePath)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}
*/

func GetTestWebServer(port int) (webServer Server, err error) {
	toReturn := NewTestWebServer()

	err = toReturn.SetPort(port)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

func MustGetTestWebServer(port int) (webServer Server) {
	webServer, err := GetTestWebServer(port)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return webServer
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

func (t *TestWebServer) MustGetMux() (mux *http.ServeMux) {
	mux, err := t.GetMux()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return mux
}

func (t *TestWebServer) MustGetPort() (port int) {
	port, err := t.GetPort()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return port
}

func (t *TestWebServer) MustGetServer() (server *http.Server) {
	server, err := t.GetServer()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return server
}

func (t *TestWebServer) MustGetWebServerWaitGroup() (webServerWaitGroup *sync.WaitGroup) {
	webServerWaitGroup, err := t.GetWebServerWaitGroup()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return webServerWaitGroup
}

func (t *TestWebServer) MustSetMux(mux *http.ServeMux) {
	err := t.SetMux(mux)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (t *TestWebServer) MustSetPort(port int) {
	err := t.SetPort(port)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (t *TestWebServer) MustSetServer(server *http.Server) {
	err := t.SetServer(server)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (t *TestWebServer) MustSetWebServerWaitGroup(webServerWaitGroup *sync.WaitGroup) {
	err := t.SetWebServerWaitGroup(webServerWaitGroup)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (t *TestWebServer) MustStartInBackground(verbose bool) {
	err := t.StartInBackground(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (t *TestWebServer) MustStop(verbose bool) {
	err := t.Stop(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
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

func (t *TestWebServer) StartInBackground(verbose bool) (err error) {
	port, err := t.GetPort()
	if err != nil {
		return err
	}

	if verbose {
		logging.LogInfof(
			"Start testWebServer in background on port %d started.",
			port,
		)
	}

	if t.webServerWaitGroup == nil {
		t.webServerWaitGroup = new(sync.WaitGroup)
	} else {
		return tracederrors.TracedError(ErrWebServerAlreadyRunning)
	}

	t.mux = http.NewServeMux()
	t.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "TestWebServer main page\n")
	})

	t.mux.HandleFunc("/hello_world.txt", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello world\n")
	})

	t.mux.HandleFunc("/example1.yaml", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "---\nhello: world\n")
	})

	t.server = &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: t.mux,
	}

	t.webServerWaitGroup.Add(1)
	go func() {
		defer t.webServerWaitGroup.Done()

		if err := t.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	if verbose {
		logging.LogInfof(
			"Start testWebServer in background on port %d finished.",
			port,
		)
	}

	return nil
}

func (t *TestWebServer) Stop(verbose bool) (err error) {
	if verbose {
		logging.LogInfo(
			"Stop TestWebServer started.",
		)
	}

	if t.webServerWaitGroup == nil {
		if verbose {
			logging.LogInfof("TestWebServer already stopped")
		}
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

	if verbose {
		logging.LogInfo(
			"Stop TestWebServer finished.",
		)
	}

	return nil
}
