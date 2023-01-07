package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/metinorak/wspubsubserver/handler"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var port = "8080"

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	godotenv.Load()

	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logWriter := zerolog.SyncWriter(os.Stdout)
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		logWriter = zerolog.NewConsoleWriter()
	}

	logger := zerolog.New(logWriter).With().Caller().Timestamp().Logger()
	ctx = logger.WithContext(ctx)

	if err := run(ctx); err != nil {
		logger.Fatal().Stack().Err(err).Msgf("program exited with an error: %+v", err)
	}
}

func run(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)
	router := handler.NewRouter()

	logger.Info().Msgf("Server is listening on port %s", port)

	chErr := make(chan error, 1)

	go func() {
		chErr <- http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-chErr:
		return err
	}
}
