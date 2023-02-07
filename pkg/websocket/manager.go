package websocket

import (
	context_config "dvm.wallet/harsh/cmd/api/context"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	//ErrEventNotSupported = errors.New("this Event type is not supported")
)

type Manager struct {
	Clients          ClientList
	ClientUserIDList ClientIDList
	sync.Mutex
	handlers map[string]EventHandler
}

func NewManager() *Manager {
	m := &Manager{
		Clients:          make(ClientList),
		ClientUserIDList: make(ClientIDList),
		handlers:         make(map[string]EventHandler),
	}
	//m.setupEventHandlers()
	return m
}

func (m *Manager) ServeWs(w http.ResponseWriter, r *http.Request) {
	log.Println("NEW CONNECTON")
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := NewClient(conn, m)
	usr := context_config.ContextGetAuthenticatedUser(r)
	m.addClient(client, usr.ID)
	go client.readMessage(usr.ID)
	go client.writeMessage(usr.ID)
}

func (m *Manager) addClient(client *Client, clientId int) {
	m.Lock()
	defer m.Unlock()
	m.Clients[client] = true
	m.ClientUserIDList[clientId] = client
}

func (m *Manager) removeClient(client *Client, clientId int) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.Clients[client]; ok {
		client.connection.Close()
		delete(m.Clients, client)
		delete(m.ClientUserIDList, clientId)
	}
}

//
//func (m *Manager) setupEventHandlers() {
//	//m.handlers[EventUpdateBalance] = func(e Event, c *Client) error {
//	//	fmt.Println(e)
//	//	return nil
//	//}
//	m.handlers[EventUpdateBalance] = UpdateBalanceHandler
//	m.handlers[EventUpdateOrder] = UpdateOrderStatusHandler
//}
//
//func (m *Manager) routeEvent(event Event, c *Client) error {
//	if handler, ok := m.handlers[event.Type]; ok {
//		if err := handler(event, c); err != nil {
//			return err
//		}
//		return nil
//	} else {
//		return ErrEventNotSupported
//	}
//}
