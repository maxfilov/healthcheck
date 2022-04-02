package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	server        *http.Server
	errors        chan error
	logger        log.FieldLogger
	listenAddress string
}

func MakeServer(port int, handler http.Handler, logger log.FieldLogger) *Server {
	http.Handle("/health", handler)
	listenAddress := fmt.Sprintf(":%d", port)
	server := &http.Server{Addr: listenAddress}
	errors := make(chan error, 1)
	return &Server{server, errors, logger, listenAddress}
}

func (server *Server) StartAsync() {
	server.logger.Info("starting http server")
	go func() {
		server.logger.Infof("listening on '%s'", server.listenAddress)
		server.errors <- server.server.ListenAndServe()
	}()
}

func (server *Server) Shutdown() error {
	server.logger.Info("stopping http server")
	return server.server.Shutdown(context.TODO())
}

func (server *Server) AwaitShutdown() error {
	server.logger.Info("waiting for http server to stop")
	result := <-server.errors
	server.logger.Info("http server stopped")
	if result == nil || result == http.ErrServerClosed {
		return nil
	} else {
		return result
	}
}
