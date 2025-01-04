package http

import (
	"context"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/asciich/asciichgolangpublic"
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
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return webServer
}

func NewTestWebServer() (t *TestWebServer) {
	return new(TestWebServer)
}

func (t *TestWebServer) GetPort() (port int, err error) {
	if t.port <= 0 {
		return -1, asciichgolangpublic.TracedError("port not set")
	}

	return t.port, nil
}

func (t *TestWebServer) GetServer() (server *http.Server, err error) {
	if t.server == nil {
		return nil, asciichgolangpublic.TracedErrorf("server not set")
	}

	return t.server, nil
}

func (t *TestWebServer) GetWebServerWaitGroup() (webServerWaitGroup *sync.WaitGroup, err error) {
	if t.webServerWaitGroup == nil {
		return nil, asciichgolangpublic.TracedErrorf("webServerWaitGroup not set")
	}

	return t.webServerWaitGroup, nil
}

func (t *TestWebServer) MustGetPort() (port int) {
	port, err := t.GetPort()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return port
}

func (t *TestWebServer) MustGetServer() (server *http.Server) {
	server, err := t.GetServer()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return server
}

func (t *TestWebServer) MustGetWebServerWaitGroup() (webServerWaitGroup *sync.WaitGroup) {
	webServerWaitGroup, err := t.GetWebServerWaitGroup()
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}

	return webServerWaitGroup
}

func (t *TestWebServer) MustSetPort(port int) {
	err := t.SetPort(port)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (t *TestWebServer) MustSetServer(server *http.Server) {
	err := t.SetServer(server)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (t *TestWebServer) MustSetWebServerWaitGroup(webServerWaitGroup *sync.WaitGroup) {
	err := t.SetWebServerWaitGroup(webServerWaitGroup)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (t *TestWebServer) MustStartInBackground(verbose bool) {
	err := t.StartInBackground(verbose)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (t *TestWebServer) MustStop(verbose bool) {
	err := t.Stop(verbose)
	if err != nil {
		asciichgolangpublic.LogGoErrorFatal(err)
	}
}

func (t *TestWebServer) SetPort(port int) (err error) {
	if port <= 0 {
		return asciichgolangpublic.TracedErrorf("Invalid value '%d' for port", port)
	}

	t.port = port

	return nil
}

func (t *TestWebServer) SetServer(server *http.Server) (err error) {
	if server == nil {
		return asciichgolangpublic.TracedErrorf("server is nil")
	}

	t.server = server

	return nil
}

func (t *TestWebServer) SetWebServerWaitGroup(webServerWaitGroup *sync.WaitGroup) (err error) {
	if webServerWaitGroup == nil {
		return asciichgolangpublic.TracedErrorf("webServerWaitGroup is nil")
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
		asciichgolangpublic.LogInfof(
			"Start testWebServer in background on port %d started.",
			port,
		)
	}

	if t.webServerWaitGroup == nil {
		t.webServerWaitGroup = new(sync.WaitGroup)
	} else {
		return asciichgolangpublic.TracedError(ErrWebServerAlreadyRunning)
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
		asciichgolangpublic.LogInfof(
			"Start testWebServer in background on port %d finished.",
			port,
		)
	}

	return nil
}

func (t *TestWebServer) Stop(verbose bool) (err error) {
	if verbose {
		asciichgolangpublic.LogInfo(
			"Stop TestWebServer started.",
		)
	}

	if t.webServerWaitGroup == nil {
		if verbose {
			asciichgolangpublic.LogInfof("TestWebServer already stopped")
		}
		return nil
	}

	if t.server == nil {
		return asciichgolangpublic.TracedError("Unexpected t.server == nil")
	}

	err = t.server.Shutdown(context.TODO())
	if err != nil {
		return asciichgolangpublic.TracedErrorf(
			"Shutdown TestWebServer failed: '%w'",
			err,
		)
	}

	t.webServerWaitGroup.Wait()
	t.webServerWaitGroup = nil

	if verbose {
		asciichgolangpublic.LogInfo(
			"Stop TestWebServer finished.",
		)
	}

	return nil
}
