package webserver

import (
	"go.uber.org/zap"
)

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

type Server struct {
	// cfg    *config.Config
	logger *zap.Logger
}

func NewServer(logger *zap.Logger) *Server {
	return &Server{
		// cfg:    cfg,
		logger: logger}
}
