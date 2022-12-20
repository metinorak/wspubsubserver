package listener

import (
	"context"
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/metinorak/wspubsub"
	"github.com/metinorak/wspubsubserver/entity"
	"github.com/rs/zerolog"
)

func ListenConnection(ctx context.Context, conn *websocket.Conn, pubSubManager *wspubsub.VartoWS) {
	logger := zerolog.Ctx(ctx)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			logger.Error().Err(err).Msg("Error while reading message")
			break
		}

		var payload entity.WsPayload
		err = json.Unmarshal(message, &payload)
		if err != nil {
			logger.Error().Err(err).Msg("Error while unmarshalling message")
			continue
		}

		switch payload.Action {
		case "SUBSCRIBE":
			pubSubManager.Subscribe(conn, payload.Topic)
		case "UNSUBSCRIBE":
			pubSubManager.Unsubscribe(conn, payload.Topic)
		case "PUBLISH":
			pubSubManager.Publish(payload.Topic, []byte(payload.Message))
		case "BROADCASTALL":
			pubSubManager.BroadcastToAll([]byte(payload.Message))
		}
	}
}
