package listener

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/metinorak/varto"
	"github.com/metinorak/wspubsubserver/entity"
	"github.com/rs/zerolog"
)

func ListenConnection(ctx context.Context, conn varto.Connection, pubSubManager *varto.Varto) {
	logger := zerolog.Ctx(ctx)
	defer func() {
		err := pubSubManager.RemoveConnection(conn)
		if err != nil {
			logger.Error().Err(err).Msg("Error while removing connection")
		}
	}()

	for {
		message, err := conn.Read()
		if err != nil {
			logger.Error().Err(err).Msg("Error while reading message")
			break
		}

		var payload entity.WsPayload
		err = json.Unmarshal(message, &payload)
		if err != nil {
			handleInvalidMessage(conn, message)
			continue
		}

		// make action case insensitive
		payload.Action = strings.ToUpper(payload.Action)

		switch payload.Action {
		case "SUBSCRIBE":
			handleSubscribe(conn, pubSubManager, &payload)
		case "UNSUBSCRIBE":
			handleUnsubscribe(conn, pubSubManager, &payload)
		case "PUBLISH":
			handlePublish(conn, pubSubManager, &payload)
		case "BROADCASTALL":
			handleBroadcastToAll(conn, pubSubManager, &payload)
		default:
			handleUnknown(conn, &payload)
		}
	}
}
