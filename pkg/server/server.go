package server

import (
	"os"
	"strings"

	"github.com/panjf2000/gnet/v2"
	"golang.org/x/exp/slog"
)

type agentCheckServer struct {
	*gnet.BuiltinEventEngine
	logger *slog.Logger
}

func (s *agentCheckServer) OnBoot(eng gnet.Engine) (action gnet.Action) {
	s.logger.Info("Health check server is running")
	return
}

func (s *agentCheckServer) OnTraffic(c gnet.Conn) gnet.Action {
	// Read the incoming message
	buffer, _ := c.Next(-1)
	message := string(buffer)

	// Log the received message
	s.logger.Info("Received message", "message", message)

	// Respond with "OK" if the message is "ping", otherwise respond with "ERROR"
	response := "ERROR"
	if strings.TrimSpace(message) == "ping" {
		response = "OK"
	}

	// Define the callback function
	callback := func(c gnet.Conn, err error) error {
		if err != nil {
			s.logger.Error("Failed to send response", "error", err)
		} else {
			s.logger.Info("Response sent", "response", response)
		}
		// Close the connection after responding
		c.Close()
		return err
	}

	// Use AsyncWrite with a callback function
	err := c.AsyncWrite([]byte(response), callback)
	if err != nil {
		s.logger.Error("Failed to send response", "error", err)
	}

	return gnet.None
}

func Run() {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler)

	// Create a new agentCheckServer
	server := &agentCheckServer{
		logger: logger,
	}

	// Start the gnet server
	err := gnet.Run(server, "tcp://:9000", gnet.WithMulticore(true))
	if err != nil {
		logger.Error("Failed to start server", "error", err)
	}
	return
}
