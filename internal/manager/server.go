package manager

import (
	"errors"
	"sync"
)

const ServerMaxClientNum = 1000 //群员最大数量

// Server server
type Server struct {
	address    string
	clientList sync.Map
	isRunning  bool
}

func (s *Server) GetAddress() string {
	return s.address
}

func (s *Server) AddClient(client *Client) {
	oldValue, ok := s.clientList.Load(client.GetClientId())
	if ok {
		oldClient := oldValue.(*Client)
		oldClient.Close(errors.New("kick off"))
	}
	s.clientList.Store(client.GetClientId(), client)
}

func (s *Server) RemoveClient(client *Client) {
	s.clientList.Delete(client.GetClientId())
}

func (s *Server) GetClientList() []string {
	var list []string
	s.clientList.Range(func(key, value interface{}) bool {
		item := value.(string)
		list = append(list, item)
		return true
	})
	return list
}

func (s *Server) CountClient() int {
	list := s.GetClientList()
	return len(list)
}

func (s *Server) Close() {
	s.clientList.Range(func(key, value interface{}) bool {
		s.clientList.Delete(key)
		return true
	})
	s.isRunning = false
}

func NewServer(address string) *Server {
	server := &Server{
		address:    address,
		clientList: sync.Map{},
		isRunning:  true,
	}
	return server
}
