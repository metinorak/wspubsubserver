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

	for {
		message, err := conn.Read()
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

		// make action case insensitive
		action := strings.ToUpper(payload.Action)

		switch action {
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
