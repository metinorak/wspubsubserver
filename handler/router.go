package handler

import "net/http"

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", HandleWebSocket)

	return mux
}
