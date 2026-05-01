package http

import (
	"context"
	"net/http"
	"time"
)

type (
	SseServer struct {
		ctx  context.Context
		w    http.ResponseWriter
		data chan string
	}
	SseResponseWriter struct {
		http.ResponseWriter
	}
)

var (
	preID      = "id: "
	preData    = "data: "
	preEvent   = "event: "
	preRetry   = "retry: "
	preComment = ":"
	pingMsg    = ": ping"
	splitMsg   = "\n\n"

	DoneData = "[DONE]"
)

func RegisterSSE(ctx context.Context, w http.ResponseWriter, dataChan chan []byte, doneFunc func()) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	WriteFlush(w, []byte(pingMsg+splitMsg))

	timerDuration := time.Second * 15
	timer := time.NewTimer(timerDuration)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			if doneFunc != nil {
				doneFunc()
			}
			return
		case data := <-dataChan:
			WriteFlush(w, []byte(preData+string(data)+splitMsg))
		case <-timer.C:
			WriteFlush(w, []byte(pingMsg+splitMsg))
			timer.Reset(timerDuration)
		}
	}
}

func WriteFlush(w http.ResponseWriter, msg []byte) {
	_, _ = w.Write(msg)
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}
