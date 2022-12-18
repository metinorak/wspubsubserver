package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/metinorak/wspubsub"
	"github.com/metinorak/wspubsubserver/listener"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var logger zerolog.Logger

// env variables
var (
	port string
)

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	if port == "" {
		port = "8080"
	}

	if err := godotenv.Load(); err != nil {
		logger.Error().Err(err).Msg("Error while loading .env file")
	}

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
}

func main() {
	ctx := logger.WithContext(context.Background())

	if err := run(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Server is shutting down")
	}
}

func run(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)

	// Create a new instance of the server
	pubSubManager := wspubsub.New()

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		conn.SetCloseHandler(func(code int, text string) error {
			return pubSubManager.RemoveConnection(conn)
		})

		pubSubManager.AddConnection(conn)

		// Listen the connection
		go listener.ListenConnection(ctx, conn, pubSubManager)
	})

	logger.Info().Msgf("Server is listening on port %s", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
		logger.Error().Err(err).Msg("Server is shutting down")
	}

	return nil
}