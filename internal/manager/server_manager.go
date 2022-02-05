package manager

import "sync"

type ServerManager struct {
	serverList sync.Map
}

func (m *ServerManager) Add(server *Server) {
	m.serverList.Store(server.GetAddress(), server)
}

func (m *ServerManager) Remove(server *Server) {
	m.serverList.Delete(server.GetAddress())
}

func (m *ServerManager) GetServerList() []string {
	var list []string
	m.serverList.Range(func(key, value interface{}) bool {
		item := value.(string)
		list = append(list, item)
		return true
	})
	return list
}

func (m *ServerManager) GetServerByAddress(address string) *Server {
	server, ok := m.serverList.Load(address)
	if !ok {
		return nil
	}
	return server.(*Server)
}

var (
	onceServerManager     sync.Once
	instanceServerManager *ServerManager
)

func GetServerManagerInstance() *ServerManager {
	onceServerManager.Do(func() {
		instanceServerManager = &ServerManager{}
	})
	return instanceServerManager
}
