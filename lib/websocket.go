package lib

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var (
	once sync.Once
	Upgrade websocket.Upgrader
)

func init() {
	once.Do(func() {
		Upgrade = websocket.Upgrader{
			ReadBufferSize: 1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
	})
}
