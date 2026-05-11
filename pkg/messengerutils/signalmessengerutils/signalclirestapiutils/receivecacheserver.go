package signalclirestapiutils

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/datetime/durationparser"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/messengerutils/signalmessengerutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type ReceiveCacheServer struct {
	mu          sync.Mutex
	cache       []*signalmessengerutils.Message
	cacheSize   int
	receiveChan chan *signalmessengerutils.Message
	stopChan    chan struct{}
}

func NewReceiveCacheServer(ctx context.Context, options *ReceiveCacheServerOptions) (*ReceiveCacheServer, error) {
	if options.CacheSize <= 0 {
		return nil, tracederrors.TracedErrorf("Invalid cache size: %d", options.CacheSize)
	}

	server := &ReceiveCacheServer{
		cache:       make([]*signalmessengerutils.Message, 0, options.CacheSize),
		cacheSize:   options.CacheSize,
		receiveChan: make(chan *signalmessengerutils.Message),
		stopChan:    make(chan struct{}),
	}

	intervalDuration, err := durationparser.ToSecondsAsTimeDuration(options.Interval)
	if err != nil {
		return nil, err
	}

	accountNumber, err := options.GetAccountNumber()
	if err != nil {
		return nil, err
	}

	go server.startReceiver(ctx, options.SignalResetClientApiUrl, *intervalDuration, accountNumber)
	return server, nil

}

func (s *ReceiveCacheServer) startReceiver(ctx context.Context, apiUrl string, interval time.Duration, accountNumber string) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var receiveFunc = func() {
		messages, err := ReceiveMessages(ctx, apiUrl, accountNumber)
		if err != nil {
			logging.LogErrorf("Failed to receive messages: %v", err)
			return
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		for _, msg := range messages {
			s.cache = append(s.cache, msg)
			if len(s.cache) > s.cacheSize {
				s.cache = s.cache[:len(s.cache)-1]
			}
		}

		logging.LogInfoByCtxf(ctx, "There are currently %d received signal messages in the cache.", len(s.cache))
	}

	receiveFunc()

	for {
		select {
		case <-ticker.C:
			receiveFunc()
		case <-s.stopChan:
			return
		}
	}
}

func (s *ReceiveCacheServer) GetMessages() []*signalmessengerutils.Message {
	s.mu.Lock()
	defer s.mu.Unlock()

	ret := []*signalmessengerutils.Message{}
	ret = append(ret, s.cache...)
	return ret
}

func (s *ReceiveCacheServer) Stop() {
	close(s.stopChan)
}

func RunReceiveCacheServer(ctx context.Context, options *ReceiveCacheServerOptions) error {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	logging.LogInfoByCtx(ctx, "Run receive cache server for signal cli rest api started.")

	server, err := NewReceiveCacheServer(ctx, options)
	if err != nil {
		return err
	}
	defer server.Stop()

	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		messages := server.GetMessages()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(messages)
	})

	errChan := make(chan error)
	go func() {
		const port = 8080
		logging.LogInfoByCtxf(ctx, "Going to start signal receive cache server on port %d.", port)
		err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
		if err != nil && err != http.ErrServerClosed {
			errChan <- tracederrors.TracedErrorf("Failed to listen and serve: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		logging.LogErrorf("Failed to start server: %v", err)
		return err
	}
}
