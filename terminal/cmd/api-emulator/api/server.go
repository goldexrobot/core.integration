package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"sync"

	apiv1 "github.com/goldexrobot/core.integration/terminal/api/v1"
	"github.com/goldexrobot/core.integration/terminal/cmd/api-emulator/api/jsonrpc"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

type Server struct {
	logger *logrus.Entry
	lis    net.Listener
	rs     *rpc.Server
}

func NewServer(port int, a *apiv1.Impl, logger *logrus.Entry) (*Server, error) {
	rs := rpc.NewServer()
	if err := rs.Register(&RPC{a}); err != nil {
		return nil, err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		logger: logger,
		lis:    lis,
		rs:     rs,
	}
	return s, nil
}

func (s *Server) Serve(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := sync.WaitGroup{}

	// ws server
	ws := &websocket.Server{
		Handler:   wsConnectionHandler(s.rs),
		Handshake: func(c *websocket.Config, r *http.Request) error { return nil },
	}

	// http server
	http.Handle("/ws", ws)
	hs := &http.Server{
		Handler: http.DefaultServeMux,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
		if err := hs.Serve(s.lis); err != nil {
			if err != http.ErrServerClosed {
				s.logger.WithError(err).Errorf("Failed to start HTTP server")
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
		<-ctx.Done()
		hs.Close()
	}()

	s.logger.Infof("Websocket handler at %v/ws", s.lis.Addr())

	wg.Wait()
}

func wsConnectionHandler(rs *rpc.Server) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		codec := jsonrpc.NewServerCodec(c, "RPC")
		rs.ServeCodec(codec)
	}
}
