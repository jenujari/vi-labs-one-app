package main

import (
	"errors"
	"fmt"
	"headless/config"
	"headless/handler"
	"net/http"
)

func main() {
	defer func() {
		close(handler.ZerodhaRequestToken)
	}()

	config := config.GetConfig()

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", config.App.Port),
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		MaxHeaderBytes:    0,
	}

	r := http.NewServeMux()

	handler.SetRoutes(r)

	server.Handler = r

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		panic(err) // failure/timeout shutting down the server gracefully
	}
}
