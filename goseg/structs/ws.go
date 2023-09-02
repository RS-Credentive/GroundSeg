package structs
// contained herein: structs for managing mutexed maps of
// mutexed websocket connections to avoid panics;
// actual writing is done via broadcast package

import (
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil)) // for debug
)

// wrapped ws+mutex
type MuConn struct {
	Conn *websocket.Conn
	Mu   sync.RWMutex
}

// mutexed ws write
func (ws *MuConn) Write(data []byte) error {
	ws.Mu.Lock()
	logger.Info("conn mu lock")
	defer ws.Mu.Unlock()
	return ws.Conn.WriteMessage(websocket.TextMessage, data)
}

// wrappers for mutexed token:websocket maps
// the maps are also mutexed as wholes
type ClientManager struct {
	AuthClients 		 map[string]*MuConn
	UnauthClients        map[string]*MuConn
	Mu                   sync.RWMutex
}

// register a new connection
func (cm *ClientManager) NewConnection(conn *websocket.Conn, tokenId string) *MuConn {
	muConn := &MuConn{Conn: conn}
	cm.Mu.Lock()
	defer cm.Mu.Unlock()
	cm.UnauthClients[tokenId] = muConn
	return muConn
}

func (cm *ClientManager) AddAuthClient(id string, client *MuConn) {
	cm.Mu.Lock()
	defer cm.Mu.Unlock()
	// also remove from other map
	// delete by value rather than key
	if _, ok := cm.UnauthClients[id]; ok {
		delete(cm.UnauthClients, id)
		cm.AuthClients[id] = client
		for token, con := range cm.UnauthClients {
			if con.Conn == client.Conn {
				delete(cm.UnauthClients, token)
				logger.Info(fmt.Sprintf("removing unauth: %v",token))
			}
		}
	}
	logger.Info(fmt.Sprintf("auth map: %v",cm.AuthClients))
}

func (cm *ClientManager) AddUnauthClient(id string, client *MuConn) {
	cm.Mu.Lock()
	defer cm.Mu.Unlock()
	cm.UnauthClients[id] = client
	// also remove from other map
	if _, ok := cm.AuthClients[id]; ok {
		delete(cm.AuthClients, id)
		for token, con := range cm.AuthClients {
			if con.Conn == client.Conn {
				logger.Info(fmt.Sprintf("removing auth: %v",token))
				delete(cm.AuthClients, token)
			}
		}
	}
	logger.Info(fmt.Sprintf("unauth map: %v",cm.AuthClients))
}

func (cm *ClientManager) BroadcastUnauth(data []byte) {
	cm.Mu.RLock()
	defer cm.Mu.RUnlock()
	logger.Info(fmt.Sprintf("Locking for %v unauth clients",len(cm.UnauthClients)))
	for _, client := range cm.UnauthClients {
		logger.Info("auth broadcasting")
		client.Write(data)
	}
	logger.Info("Unlocking")
}

func (cm *ClientManager) BroadcastAuth(data []byte) {
	cm.Mu.RLock()
	defer cm.Mu.RUnlock()
	logger.Info(fmt.Sprintf("Locking for %v auth clients",len(cm.AuthClients)))
	for _, client := range cm.AuthClients {
		logger.Info("unauth broadcasting")
		client.Write(data)
	}
	logger.Info("Unlocking")
}

type WsType struct {
	Payload struct {
		Type string `json:"type"`
	} `json:"payload"`
}

type WsPayload struct {
	ID      string        `json:"id"`
	Type    string        `json:"type"`
	Payload interface{}   `json:"payload"`
	Token   WsTokenStruct `json:"token"`
}

type WsUrbitPayload struct {
	ID      string        `json:"id"`
	Type    string        `json:"type"`
	Payload WsUrbitAction `json:"payload"`
	Token   WsTokenStruct `json:"token"`
}

type WsSystemPayload struct {
	ID      string         `json:"id"`
	Type    string         `json:"type"`
	Payload WsSystemAction `json:"payload"`
	Token   WsTokenStruct  `json:"token"`
}

type WsUrbitAction struct {
	Type   string `json:"type"`
	Action string `json:"action"`
	Patp   string `json:"patp"`
}

type WsSystemAction struct {
	Type    string `json:"type"`
	Action  string `json:"action"`
	Command string `json:"command"`
}

type WsTokenStruct struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

type WsLoginPayload struct {
	ID      string        `json:"id"`
	Type    string        `json:"type"`
	Payload WsLoginAction `json:"payload"`
	Token   WsTokenStruct `json:"token"`
}

type WsLoginAction struct {
	Type     string `json:"type"`
	Password string `json:"password"`
}

type WsPwPayload struct {
	ID      string `json:"id"`
	Payload WsPwAction `json:"payload"`
	Token    WsTokenStruct `json:"token"`
}

type WsPwAction struct {
	Type     string `json:"type"`
	Action   string `json:"action"`
	Old      string `json:"old"`
	Password string `json:"password"`
}

type WsLogoutPayload struct {
	ID      string        `json:"id"`
	Token   WsTokenStruct `json:"token"`
}

type WsResponsePayload struct {
	ID       string        `json:"id"`
	Type     string        `json:"type"`
	Response string        `json:"response"`
	Error    string        `json:"error"`
	Token    WsTokenStruct `json:"token"`
}

type WsStartramPayload struct {
	ID      string        `json:"id"`
	Type    string        `json:"type"`
	Payload WsStartramAction `json:"payload"`
	Token   WsTokenStruct `json:"token"`
}

type WsStartramAction struct {
	Type   string `json:"type"`
	Action string `json:"action"`
	Key    string `json:"key"`
	Region string `json:"region"`
}