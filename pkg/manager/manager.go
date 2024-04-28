package manager

import (
	"context"
	"errors"
	"sync"

	"github.com/Hank-Kuo/chat-app/config"

	"github.com/redis/go-redis/v9"
)

type ClientManager struct {
	Ctx             context.Context
	cfg             *config.Config
	InstanceId      string
	ClientIdMap     map[string]*Client
	ClientIdMapLock sync.RWMutex
	Connect         chan *Client
	DisConnect      chan *Client
	Rdb             *redis.Client
	ToClientChan    chan ToClientInfo
}

func NewClientManager(rdb *redis.Client, cfg *config.Config) *ClientManager {
	return &ClientManager{
		Ctx:          context.Background(),
		cfg:          cfg,
		InstanceId:   cfg.Server.InstanceIP,
		ClientIdMap:  make(map[string]*Client),
		Connect:      make(chan *Client, 10000),
		DisConnect:   make(chan *Client, 10000),
		Rdb:          rdb,
		ToClientChan: make(chan ToClientInfo, 1000),
	}
}

func (manager *ClientManager) Start() {
	for {
		select {
		case client := <-manager.Connect:
			manager.AddClient(client)
		case client := <-manager.DisConnect:
			_ = client.Socket.Close()
			manager.DelClient(client)
			client.IsDeleted = true
			client = nil
		}
	}
}

func (m *ClientManager) AddClient(client *Client) error {
	m.ClientIdMapLock.Lock()
	defer m.ClientIdMapLock.Unlock()

	err := m.Rdb.Set(m.Ctx, client.ClientId, m.InstanceId, 0).Err()
	if err != nil {
		return err
	}

	// m.Rdb.Set(m.Ctx, "bac09a89-df1a-4644-ba2f-89f4da8d0456", "localhost2", 0)
	m.ClientIdMap[client.ClientId] = client
	return nil
}

func (m *ClientManager) DelClient(client *Client) error {

	m.ClientIdMapLock.Lock()
	defer m.ClientIdMapLock.Unlock()

	delete(m.ClientIdMap, client.ClientId)

	_, err := m.Rdb.Del(m.Ctx, client.ClientId).Result()
	if err != nil {
		return err
	}
	return nil
}

func (m *ClientManager) Count() int {
	m.ClientIdMapLock.RLock()
	defer m.ClientIdMapLock.RUnlock()
	return len(m.ClientIdMap)
}

func (m *ClientManager) GetByClientId(clientId string) (*Client, error) {
	m.ClientIdMapLock.RLock()
	defer m.ClientIdMapLock.RUnlock()

	if client, ok := m.ClientIdMap[clientId]; !ok {
		return nil, errors.New("client not exist")
	} else {
		return client, nil
	}
}

func (m *ClientManager) GetInstacesByClients(clientId string) (string, error) {
	return m.Rdb.Get(m.Ctx, clientId).Result()
}
