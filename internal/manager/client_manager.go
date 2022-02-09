package manager

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"sync"
)

type ClientManager struct {
	clientUserIdMap sync.Map
	clientIdClientMap sync.Map
	connCh          chan *Client
	disConnCh       chan *Client
	broadcastCh     chan *Client
	once  sync.Once
	mutex sync.Mutex
}

func (m *ClientManager) Connect(client *Client) {
	if client.clientId!="" && client.userId!=0{
		m.connCh <- client
		m.broadcastCh <- client
	}
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
			color.Red("client connect success:%s", client.clientId)
			m.clientUserIdMap.Store(client, client.userId)
			m.clientIdClientMap.Store(client.clientId,client)
			break
		case client := <-m.disConnCh:
			m.clientUserIdMap.Delete(client)
			m.clientIdClientMap.Delete(client.clientId)
			break
		case _ = <-m.broadcastCh:
			total := m.GetUserIdCnt()
			m.clientUserIdMap.Range(func(key, value interface{}) bool {
				_client := key.(*Client)
				m := fmt.Sprintf("user %d logined,client total:%d", _client.userId, total)
				_client.Send([]byte(m))
				return true
			})
			break
		}
	}
}
func (m *ClientManager) Close() {
	m.mutex.Lock()
	close(m.connCh)
	close(m.disConnCh)
	m.clientUserIdMap.Range(func(key, value interface{}) bool {
		client := value.(*Client)
		client.Close(errors.New("client manager closed"))
		return true
	})
	m.mutex.Unlock()
}

func (m *ClientManager) GetClientByClientId(clientId string) *Client  {
	value,ok:= m.clientIdClientMap.Load(clientId)
	if !ok{
		return nil
	}
	return value.(*Client)
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
			mutex:           sync.Mutex{},
		}
	})
	return instanceClientManager
}
