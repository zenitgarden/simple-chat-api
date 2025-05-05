package logger

import (
	"log"
	"time"

	"github.com/zenitgarden/simple-chat-api/internal/handler"
)

func StartConnectionLogger() {
	go func() {
		for {
			log.Printf("Active WS connections: %d", handler.ConnectionCount())
			time.Sleep(10 * time.Second)
		}
	}()
}
