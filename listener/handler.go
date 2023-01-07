package listener

import (
	"encoding/json"
	"fmt"

	"github.com/metinorak/varto"
	"github.com/metinorak/wspubsubserver/entity"
)

func handleResponseMessage(conn varto.Connection, response *entity.WsResponse) {
	bs, err := json.Marshal(response)
	if err != nil {
		conn.Write([]byte(err.Error()))
	} else {
		conn.Write(bs)
	}
}

func handleSubscribe(conn varto.Connection, pubSubManager *varto.Varto, payload *entity.WsPayload) {
	err := pubSubManager.Subscribe(conn, payload.Topic)
	if err != nil {
		handleResponseMessage(conn, &entity.WsResponse{
			Action:  payload.Action,
			Topic:   payload.Topic,
			Message: err.Error(),
			Status:  "ERROR",
		})
	} else {
		handleResponseMessage(conn, &entity.WsResponse{
			Action: payload.Action,
			Topic:  payload.Topic,
			Status: "OK",
		})
	}
}

func handleUnsubscribe(conn varto.Connection, pubSubManager *varto.Varto, payload *entity.WsPayload) {
	err := pubSubManager.Unsubscribe(conn, payload.Topic)
	if err != nil {
		handleResponseMessage(conn, &entity.WsResponse{
			Action:  payload.Action,
			Topic:   payload.Topic,
			Message: err.Error(),
			Status:  "ERROR",
		})
	} else {
		handleResponseMessage(conn, &entity.WsResponse{
			Action: payload.Action,
			Topic:  payload.Topic,
			Status: "OK",
		})
	}
}

func handlePublish(conn varto.Connection, pubSubManager *varto.Varto, payload *entity.WsPayload) {
	err := pubSubManager.Publish(payload.Topic, []byte(payload.Message))
	if err != nil {
		handleResponseMessage(conn, &entity.WsResponse{
			Action:  payload.Action,
			Topic:   payload.Topic,
			Message: err.Error(),
			Status:  "ERROR",
		})
	} else {
		handleResponseMessage(conn, &entity.WsResponse{
			Action: payload.Action,
			Topic:  payload.Topic,
			Status: "OK",
		})
	}
}

func handleBroadcastToAll(conn varto.Connection, pubSubManager *varto.Varto, payload *entity.WsPayload) {
	err := pubSubManager.BroadcastToAll([]byte(payload.Message))
	if err != nil {
		handleResponseMessage(conn, &entity.WsResponse{
			Action:  payload.Action,
			Message: err.Error(),
			Status:  "ERROR",
		})
	} else {
		handleResponseMessage(conn, &entity.WsResponse{
			Action: payload.Action,
			Status: "OK",
		})
	}
}

func handleUnknown(conn varto.Connection, payload *entity.WsPayload) {
	handleResponseMessage(conn, &entity.WsResponse{
		Action:  payload.Action,
		Message: "Action is not valid",
		Status:  "ERROR",
	})
}

func handleInvalidMessage(conn varto.Connection, message []byte) {
	handleResponseMessage(conn, &entity.WsResponse{
		Action:  "UNKNOWN",
		Message: fmt.Sprintf("Invalid message: %s", string(message)),
		Status:  "ERROR",
	})
}
