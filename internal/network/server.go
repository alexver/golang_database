package network

import (
	"fmt"
	"io"
	"log"
	"net"

	database "github.com/alexver/golang_database/internal"
	"github.com/alexver/golang_database/internal/config"
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.Logger
	db     *database.Database

	network string
	address string
}

func CreateServer(config *config.NetworkServerConfig, db *database.Database, logger *zap.Logger) *Server {
	if logger == nil {
		log.Fatal("Logger is not valid")
	}

	if db == nil {
		log.Fatal("Database is not valid")
	}

	return &Server{
		network: config.Network,
		address: fmt.Sprintf("%s:%d", config.Host, config.Port),
		db:      db,
		logger:  logger,
	}
}

func (s *Server) StartServer() {
	listener, err := net.Listen(s.network, s.address)
	if err != nil {
		s.logger.Fatal(fmt.Sprintf("Start TCP Server error: %s", err))
	}

	defer listener.Close()

	s.logger.Info(fmt.Sprintf("TCP Server has been started on address '%s'", s.address))

	for {
		connection, err := listener.Accept()
		if err != nil {
			s.logger.Error(fmt.Sprintf("Cannot accept connection. Error: %s", err))

			continue
		}

		go s.handleClient(connection)
	}
}

func (s *Server) handleClient(connection net.Conn) {

	defer connection.Close()

	buffer := make([]byte, 1024)

	end, err := connection.Read(buffer)
	if err != nil && err != io.EOF {
		s.logger.Error(fmt.Sprintf("Connection reader error: %s", err))
		return
	}

	s.logger.Info(fmt.Sprintf("Server got message '%s'", buffer[:end]))

	result, err := s.db.ProcessQuery(string(buffer[:end]))
	var str string
	if err != nil {
		s.logger.Error(fmt.Sprintf("Database processing error: %s", err))
		str = fmt.Sprintf("[error] %s", err)
	} else {
		str = result.(string)
	}

	_, err = connection.Write([]byte(str))
	if err != nil {
		s.logger.Error(fmt.Sprintf("Connection write error: %s", err))
		return
	}

	s.logger.Info(fmt.Sprintf("Server send back message '%s'", str))
}
