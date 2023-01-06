package entity

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// VartoConnectionAdapter is a wrapper for websocket connection that implements varto.Connection interface
type VartoConnectionAdapter struct {
	Connection *websocket.Conn
}

func (v *VartoConnectionAdapter) Write(data []byte) error {
	return v.Connection.WriteMessage(websocket.TextMessage, data)
}

func (v *VartoConnectionAdapter) Read() ([]byte, error) {
	_, data, err := v.Connection.ReadMessage()
	return data, err
}

func (v *VartoConnectionAdapter) GetId() string {
	return fmt.Sprintf("%p", v.Connection)
}
