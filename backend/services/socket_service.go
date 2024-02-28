package services

import (
	"sync"

	"github.com/gorilla/websocket"
)

type SocketService struct {
	mu          sync.Mutex
	connections map[string]*websocket.Conn
}

func NewSocketService() *SocketService {
	return &SocketService{
		connections: make(map[string]*websocket.Conn),
	}
}

// SetConnection associates a player with a WebSocket connection.
func (s *SocketService) SetConnection(playerId string, conn *websocket.Conn) {
	// s.mu.Lock()
	// defer s.mu.Unlock()
	s.connections[playerId] = conn
	return
}

// GetConnection retrieves the WebSocket connection associated with a player.
func (s *SocketService) GetConnection(playerId string) *websocket.Conn {
	s.mu.Lock()
	defer s.mu.Unlock()
	conn := s.connections[playerId]
	return conn
}

// RemoveConnection removes the WebSocket connection associated with a player.
func (s *SocketService) RemoveConnection(playerId string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connections[playerId].Close()
	delete(s.connections, playerId)
}
