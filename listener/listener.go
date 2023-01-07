package listener

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/metinorak/varto"
	"github.com/metinorak/wspubsubserver/entity"
	"github.com/rs/zerolog"
)

func handleResponseMessage(conn varto.Connection, response *entity.WsResponse) {
	bs, err := json.Marshal(response)
	if err != nil {
		conn.Write([]byte(err.Error()))
	} else {
		conn.Write(bs)
	}
}

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
			logger.Error().Err(err).Msg("Error while unmarshalling message")
			continue
		}

		// make action case insensitive
		action := strings.ToUpper(payload.Action)

		switch action {
		case "SUBSCRIBE":
			err := pubSubManager.Subscribe(conn, payload.Topic)
			if err != nil {
				handleResponseMessage(conn, &entity.WsResponse{
					Action:  action,
					Topic:   payload.Topic,
					Message: err.Error(),
					Status:  "ERROR",
				})
			} else {
				handleResponseMessage(conn, &entity.WsResponse{
					Action: action,
					Topic:  payload.Topic,
					Status: "OK",
				})
			}
		case "UNSUBSCRIBE":
			err := pubSubManager.Unsubscribe(conn, payload.Topic)
			if err != nil {
				handleResponseMessage(conn, &entity.WsResponse{
					Action:  action,
					Topic:   payload.Topic,
					Message: err.Error(),
					Status:  "ERROR",
				})
			} else {
				handleResponseMessage(conn, &entity.WsResponse{
					Action: action,
					Topic:  payload.Topic,
					Status: "OK",
				})
			}
		case "PUBLISH":
			err := pubSubManager.Publish(payload.Topic, []byte(payload.Message))
			if err != nil {
				handleResponseMessage(conn, &entity.WsResponse{
					Action:  action,
					Topic:   payload.Topic,
					Message: err.Error(),
					Status:  "ERROR",
				})
			} else {
				handleResponseMessage(conn, &entity.WsResponse{
					Action: action,
					Topic:  payload.Topic,
					Status: "OK",
				})
			}
		case "BROADCASTALL":
			err := pubSubManager.BroadcastToAll([]byte(payload.Message))
			if err != nil {
				handleResponseMessage(conn, &entity.WsResponse{
					Action:  action,
					Message: err.Error(),
					Status:  "ERROR",
				})
			} else {
				handleResponseMessage(conn, &entity.WsResponse{
					Action: action,
					Status: "OK",
				})
			}
		}
	}
}
