package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/metinorak/varto"
	"github.com/metinorak/wspubsubserver/entity"
	"github.com/metinorak/wspubsubserver/listener"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	EnableCompression: true,
}

var pubSubManager = varto.New(nil)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	vartoConn := &entity.VartoConnectionAdapter{
		Connection: conn,
	}

	conn.SetCloseHandler(func(code int, text string) error {
		return pubSubManager.RemoveConnection(vartoConn)
	})

	conn.SetCompressionLevel(9)

	pubSubManager.AddConnection(vartoConn)

	// Listen the connection
	go listener.ListenConnection(r.Context(), vartoConn, pubSubManager)
}
