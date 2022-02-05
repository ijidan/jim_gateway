package manager

import (
	"fmt"
	"github.com/fatih/color"
	"sync"
)

type ClientManager struct {
	clientUserIdMap sync.Map
	connCh          chan *Client
	disConnCh       chan *Client
	broadcastCh     chan *Client
}

func (m *ClientManager) Connect(client *Client) {
	m.connCh <- client
	m.broadcastCh <- client
}

func (m *ClientManager) DisConnect(client *Client) {
	m.disConnCh <- client
}

func (m *ClientManager) GetUserIdList() []uint64 {
	var list []uint64
	m.clientUserIdMap.Range(func(key, value interface{}) bool {
		item := value.(uint64)
		list = append(list, item)
		return true
	})
	return list
}

func (m *ClientManager) GetUserIdCnt() int {
	list := m.GetUserIdList()
	return len(list)
}

func (m *ClientManager) Loop() {
	for {
		select {
		case client := <-m.connCh:
			color.Red("client num:%d", m.GetUserIdCnt())
			m.clientUserIdMap.Store(client, client.userId)
			break
		case client := <-m.disConnCh:
			m.clientUserIdMap.Delete(client)
			break
		case client := <-m.broadcastCh:
			total := m.GetUserIdCnt()
			m.clientUserIdMap.Range(func(key, value interface{}) bool {
				_client := key.(*Client)
				if value.(uint64) != client.userId {
					m := fmt.Sprintf("user %d logined,client total:%d", _client.userId, total)
					message := NewMessage(0, _client.userId, []byte(m))
					_client.Send(message)
				}
				return true
			})
			break
		}
	}
}
func (m *ClientManager) Close() {
	close(m.connCh)
	close(m.disConnCh)
}

var (
	onceClientManager     sync.Once
	instanceClientManager *ClientManager
)

func GetClientManagerInstance() *ClientManager {
	onceClientManager.Do(func() {
		instanceClientManager = &ClientManager{
			clientUserIdMap: sync.Map{},
			connCh:          make(chan *Client, 1000),
			disConnCh:       make(chan *Client, 1000),
			broadcastCh:     make(chan *Client, 1000),
		}
	})
	return instanceClientManager
}
