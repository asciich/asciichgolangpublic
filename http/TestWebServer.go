package http

import (
	"context"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
)

// A simple webserver mostly used for testing.
type TestWebServer struct {
	webServerWaitGroup *sync.WaitGroup
	port               int
	server             *http.Server
}

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

func (t *TestWebServer) GetPort() (port int, err error) {
	if t.port <= 0 {
		return -1, errors.TracedError("port not set")
	}

	return t.port, nil
}

func (t *TestWebServer) GetServer() (server *http.Server, err error) {
	if t.server == nil {
		return nil, errors.TracedErrorf("server not set")
	}

	return t.server, nil
}

func (t *TestWebServer) GetWebServerWaitGroup() (webServerWaitGroup *sync.WaitGroup, err error) {
	if t.webServerWaitGroup == nil {
		return nil, errors.TracedErrorf("webServerWaitGroup not set")
	}

	return t.webServerWaitGroup, nil
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

func (t *TestWebServer) SetPort(port int) (err error) {
	if port <= 0 {
		return errors.TracedErrorf("Invalid value '%d' for port", port)
	}

	t.port = port

	return nil
}

func (t *TestWebServer) SetServer(server *http.Server) (err error) {
	if server == nil {
		return errors.TracedErrorf("server is nil")
	}

	t.server = server

	return nil
}

func (t *TestWebServer) SetWebServerWaitGroup(webServerWaitGroup *sync.WaitGroup) (err error) {
	if webServerWaitGroup == nil {
		return errors.TracedErrorf("webServerWaitGroup is nil")
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
		return errors.TracedError(ErrWebServerAlreadyRunning)
	}

	t.server = &http.Server{Addr: ":" + strconv.Itoa(port)}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "TestWebServer main page\n")
	})

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
		return errors.TracedError("Unexpected t.server == nil")
	}

	err = t.server.Shutdown(context.TODO())
	if err != nil {
		return errors.TracedErrorf(
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
