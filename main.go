package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(requestLogger())

	// lets make a code only change
	_ = r

	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello, world from KodaCD!")
		// spew.Dump(map[string]bool{"foo": true})
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Info().Msgf("starting server on %s", ":8080")
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Info().Msg(err.Error())
		}
	}()

	quit := make(chan os.Signal, 10)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// close(stopCh)

	log.Warn().Msg("shutting down due to signal")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Msg("forced shutdown, timeout exceeded")
		os.Exit(1)
	}

	log.Warn().Msg("shutdown complete")
}

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		values := map[string]interface{}{
			"Method": c.Request.Method,
			"URI":    c.Request.RequestURI,
			"Status": c.Writer.Status(),
		}
		log.Info().Fields(values).Msg("request")
	}
}
